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
``gcloud compute``.


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

Create some persistent disks to use for ``/var/lib/docker``. The default root partition of CoreOS
is only around 4 GB and not enough for storing Docker images and instances. The following creates 3
disks sized at 256 GB:

.. code-block:: console

    $ gcloud compute disks create cored1 cored2 cored3 --size 256GB --type pd-standard --zone us-central1-c

    NAME   ZONE          SIZE_GB TYPE        STATUS
    cored1 us-central1-c 256     pd-standard READY
    cored2 us-central1-c 256     pd-standard READY
    cored3 us-central1-c 256     pd-standard READY

Launch 3 instances. You can choose another starting CoreOS image from the listing output of
``gcloud compute images list``:

.. code-block:: console

    $ for num in 1 2 3; do
      gcloud compute instances create core${num} \
      --zone us-central1-c \
      --machine-type n1-standard-2 \
      --metadata-from-file user-data=gce-user-data,sshKeys=$HOME/.ssh/deis.pub \
      --disk name=cored${num},device-name=coredocker \
      --tags deis \
      --image coreos-stable-899-13-0-v20160323 \
      --image-project coreos-cloud;
    done

    NAME  ZONE          MACHINE_TYPE  INTERNAL_IP   EXTERNAL_IP    STATUS
    core1 us-central1-c n1-standard-2 10.240.10.107 108.59.80.10   RUNNING
    core2 us-central1-c n1-standard-2 10.240.10.108 108.59.80.11   RUNNING
    core3 us-central1-c n1-standard-2 10.240.10.109 108.59.80.12   RUNNING

.. note::

    The provision script will by default provision ``n1-standard-2`` instances. Choosing a smaller
    instance size is not recommended. Please refer to :ref:`system-requirements` for resource
    considerations when choosing an instance size to run Deis.

Load Balancing
--------------

We will need to load balance the Deis routers so we can get to Deis services (controller and builder) and our applications.

.. code-block:: console

    $ gcloud compute http-health-checks create basic-check --request-path /health-check
    $ gcloud compute target-pools create deis --health-check basic-check --session-affinity CLIENT_IP_PROTO --region us-central1
    $ gcloud compute target-pools add-instances deis --instances core1,core2,core3
    $ gcloud compute forwarding-rules create deisapp --target-pool deis --region us-central1

    NAME    REGION      IP_ADDRESS     IP_PROTOCOL TARGET
    deisapp us-central1 23.251.153.6   TCP         us-central1/targetPools/deis

Note the forwarding rule external IP address. We will use it as the Deis login endpoint in a future step. Now allow the ports on the CoreOS nodes:

.. code-block:: console

    $ gcloud compute firewall-rules create deis-router --target-tags deis --allow tcp:80,tcp:443,tcp:2222


Configure DNS
-------------

We can create DNS records in Google Cloud DNS using the ``gcloud`` utility. In our example we will
be using the domain name `deisdemo.io`. Create the zone:

.. code-block:: console

    $ gcloud dns managed-zones create --dns-name deisdemo.io. --description "Example Deis cluster domain name" deisdemoio
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

Note the `nameServers` array from the output. We will need to setup our upstream domain name
servers to these.

Now edit the zone to add the Deis endpoint and wildcard DNS:

.. code-block:: console

    $ gcloud dns record-sets transaction start --zone deisdemoio

This exports a `transaction.yaml` file.

.. code-block:: console

    ---
    additions:
    - kind: dns#resourceRecordSet
      name: deisdemo.io.
      rrdatas:
      - ns-cloud1.googledomains.com. dns-admin.google.com. 1 21600 3600 1209600 300
      ttl: 21600
      type: SOA
    deletions:
    - kind: dns#resourceRecordSet
      name: deisdemo.io.
      rrdatas:
      - ns-cloud1.googledomains.com. dns-admin.google.com. 0 21600 3600 1209600 300
      ttl: 21600
      type: SOA

You will want to add two records as YAML objects. Here is an example edit for the two A record additions:

.. code-block:: console

    ---
    additions:
    - kind: dns#resourceRecordSet
      name: deisdemo.io.
      rrdatas:
      - ns-cloud1.googledomains.com. dns-admin.google.com. 1 21600 3600 1209600 300
      ttl: 21600
      type: SOA
    - kind: dns#resourceRecordSet
      name: deis.deisdemo.io.
      rrdatas:
      - 23.251.153.6
      ttl: 21600
      type: A
    - kind: dns#resourceRecordSet
      name: "*.dev.deisdemo.io."
      rrdatas:
      - 23.251.153.6
      ttl: 21600
      type: A
    deletions:
    - kind: dns#resourceRecordSet
      name: deisdemo.io.
      rrdatas:
      - ns-cloud1.googledomains.com. dns-admin.google.com. 0 21600 3600 1209600 300
      ttl: 21600
      type: SOA

And finally execute the transaction.

.. code-block:: console

    $ gcloud dns record-sets transaction execute --zone deisdemoio


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.

It works! Enjoy your Deis cluster in Google Compute Engine!

.. _`contrib/gce`: https://github.com/deis/deis/tree/master/contrib/gce
.. _`Google Cloud SDK`: https://cloud.google.com/compute/docs/gcloud-compute/#install
.. _`Google Developer Console`: https://console.developers.google.com/project
