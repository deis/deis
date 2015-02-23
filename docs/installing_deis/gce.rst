:title: Installing Deis on Google Compute Engine
:description: How to provision a multi-node Deis cluster on Google Compute Engine

.. _deis_on_gce:

Google Compute Engine
=====================

Let's build a Deis cluster in Google's Compute Engine!

Please :ref:`get the source <get_the_source>` and refer to the scripts in `contrib/gce`_
while following this documentation.


Prerequisites
-------------

Let's get a few Google things squared away so we can provision VM instances.


Install Google Cloud SDK
^^^^^^^^^^^^^^^^^^^^^^^^

Install the `Google Cloud SDK`_. You will then need to login with your Google Account:

.. code-block:: console

    $ gcloud auth login


Create New Project
^^^^^^^^^^^^^^^^^^

Create a new project in the `Google Developer Console`_. You should get a project ID like
``orbital-gantry-285`` back. We'll set it as the default for the SDK tools:

.. code-block:: console

    $ gcloud config set project orbital-gantry-285


Enable Billing
^^^^^^^^^^^^^^

.. important::

    You will begin to accrue charges once you create resources such as disks and instances.

Navigate to the project console and then the *Billing & Settings* section in the browser. Click the
*Enable billing* button and fill out the form. This is needed to create resources in Google's
Compute Engine.


Initialize Compute Engine
^^^^^^^^^^^^^^^^^^^^^^^^^

Google Computer Engine won't be available via the command line tools until it is initialized in the
web console. Navigate to *COMPUTE* -> *COMPUTE ENGINE* -> *VM Instances* in the project console.
The Compute Engine will take a moment to initialize and then be ready to create resources via
``gcutil``.


Cloud Init
----------

Create your cloud init file using Deis' ``contrib/gce/create-gce-user-data`` script and a new etcd
discovery URL. First, install PyYAML:

.. code-block:: console

    $ sudo pip install pyyaml

Then navigate to the ``contrib/gce`` directory:

.. code-block:: console

    $ cd contrib/gce

Finally, create the ``gce-user-data`` file:

.. code-block:: console

    $ ./create-gce-user-data $(curl -s https://discovery.etcd.io/new)

We should have a ``gce-user-data`` file ready to launch CoreOS nodes with.

Launch Instances
----------------

Create a SSH key that we will use for Deis host communication:

.. code-block:: console

    $ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis

Launch 3 instances with direct-attached ephemeral ssds. You can choose another starting CoreOS image from the listing output of
``gcloud compute images list``:

.. code-block:: console

    $ for num in 1 2 3; do
      gcloud compute instances create core${num} \
      --zone us-central1-f \
      --machine-type n1-highmem-2 \
      --metadata-from-file user-data=gce-user-data sshKeys=~/.ssh/deis.pub \
      --local-ssd device-name=coredocker interface=SCSI \
      --tags deis \
      --image coreos-stable-522-6-0-v20150128 \
      --image-project coreos-cloud;
    done

    Created [https://www.googleapis.com/compute/v1/projects/deis-dialoghq/zones/us-central1-f/instances/core1].
    NAME  ZONE          MACHINE_TYPE INTERNAL_IP   EXTERNAL_IP    STATUS
    core1 us-central1-f n1-highmem-2 10.240.33.107 23.236.59.66   RUNNING
    core2 us-central1-f n1-highmem-2 10.240.94.33  108.59.80.17   RUNNING
    core3 us-central1-f n1-highmem-2 10.240.28.163 108.59.85.85   RUNNING

.. note::

    The provision script will by default provision ``n1-highmem-2`` instances. Choosing a smaller
    instance size is not recommended. Please refer to :ref:`system-requirements` for resource
    considerations when choosing an instance size to run Deis.

Load Balancing
--------------

We will need to load balance the Deis routers so we can get to Deis services (controller and builder) and our applications.

.. code-block:: console

    $ gcloud compute http-health-checks create basic-check --request-path /health-check
    $ gcloud compute target-pools create deis --health-check basic-check --session-affinity CLIENT_IP_PROTO --region us-central1
    $ gcloud compute target-pools add-instances deis --instances core1 core2 core3
    $ gcloud compute forwarding-rules create deisapp --target-pool deis --region us-central1

    Created [https://www.googleapis.com/compute/v1/projects/deis-dialoghq/regions/us-central1/forwardingRules/deisapp].
    NAME    REGION      IP_ADDRESS     IP_PROTOCOL TARGET
    deisapp us-central1 23.251.153.6   TCP         us-central1/targetPools/deis

Note the forwarding rule external IP address. We will use it as the Deis login endpoint in a future step. Now allow the ports on the CoreOS nodes:

.. code-block:: console

    $ gcloud compute firewall-rules create deis-router --target-tags deis --allow tcp:80 tcp:443 tcp:2222


Configure DNS
-------------

We can create DNS records in Google Cloud DNS using the ``gcloud`` utility. In our example we will
be using the domain name `deisdemo.io`. Create the zone:

.. code-block:: console

    $ gcloud dns managed-zone create --dns_name deisdemo.io. --description "Example Deis cluster domain name" deisdemoio
    Creating {'dnsName': 'deisdemo.io.', 'name': 'deisdemoio', 'description':
    'Example Deis cluster domain name'} in eco-theater-654

    Do you want to continue (Y/n)?  Y

    {
        "creationTime": "2014-07-28T00:01:45.835Z",
        "description": "Example Deis cluster domain name",
        "dnsName": "deisdemo.io.",
        "id": "1374035518570040348",
        "kind": "dns#managedZone",
        "name": "deisdemoio",
        "nameServers": [
            "ns-cloud-d1.googledomains.com.",
            "ns-cloud-d2.googledomains.com.",
            "ns-cloud-d3.googledomains.com.",
            "ns-cloud-d4.googledomains.com."
        ]
    }

Note the `nameServers` array from the JSON output. We will need to setup our upstream domain name
servers to these.

Now edit the zone to add the Deis endpoint and wildcard DNS:

.. code-block:: console

    $ gcloud dns records --zone deisdemoio edit
    {
        "additions": [
            {
                "kind": "dns#resourceRecordSet",
                "name": "deisdemo.io.",
                "rrdatas": [
                    "ns-cloud-d1.googledomains.com. dns-admin.google.com. 2 21600 3600 1209600 300"
                ],
                "ttl": 21600,
                "type": "SOA"
            }
        ],
        "deletions": [
            {
                "kind": "dns#resourceRecordSet",
                "name": "deisdemo.io.",
                "rrdatas": [
                    "ns-cloud-d1.googledomains.com. dns-admin.google.com. 1 21600 3600 1209600 300"
                ],
                "ttl": 21600,
                "type": "SOA"
            }
        ]
    }

You will want to add two records as JSON objects. Here is an example edit for the two A record additions:

.. code-block:: console

    {
        "additions": [
            {
                "kind": "dns#resourceRecordSet",
                "name": "deisdemo.io.",
                "rrdatas": [
                    "ns-cloud-d1.googledomains.com. dns-admin.google.com. 2 21600 3600 1209600 300"
                ],
                "ttl": 21600,
                "type": "SOA"
            },
            {
                "kind": "dns#resourceRecordSet",
                "name": "deis.deisdemo.io.",
                "rrdatas": [
                    "23.251.153.6"
                ],
                "ttl": 21600,
                "type": "A"
            },
            {
                "kind": "dns#resourceRecordSet",
                "name": "*.dev.deisdemo.io.",
                "rrdatas": [
                    "23.251.153.6"
                ],
                "ttl": 21600,
                "type": "A"
            }
        ],
        "deletions": [
            {
                "kind": "dns#resourceRecordSet",
                "name": "deisdemo.io.",
                "rrdatas": [
                    "ns-cloud-d1.googledomains.com. dns-admin.google.com. 1 21600 3600 1209600 300"
                ],
                "ttl": 21600,
                "type": "SOA"
            }
        ]
    }


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.

It works! Enjoy your Deis cluster in Google Compute Engine!

.. _`contrib/gce`: https://github.com/deis/deis/tree/master/contrib/gce
.. _`Google Cloud SDK`: https://developers.google.com/compute/docs/gcutil/#install
.. _`Google Developer Console`: https://console.developers.google.com/project
