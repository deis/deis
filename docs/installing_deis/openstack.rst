:title: Installing Deis on OpenStack
:description: How to provision a multi-node Deis cluster on OpenStack

.. _deis_on_openstack:

OpenStack
=========

Please :ref:`get the source <get_the_source>` and refer to the scripts in `contrib/openstack`_
while following this documentation.

.. important::

    OpenStack support for Deis was originally contributed by `Shlomo Hakim`_ and has been updated
    by Deis community members. OpenStack support is untested by the Deis team, so we rely on
    the community to improve this documentation and to fix bugs. We greatly appreciate the help!


Check System Requirements
-------------------------

Please refer to :ref:`system-requirements` for resource considerations when choosing a machine
size to run Deis.

Prerequisites
-------------

Make sure that the following utilities are installed and in your execution path:

* nova
* neutron
* glance

Configure OpenStack
-------------------

Create an ``openrc.sh`` file to match the following:

.. code-block:: console

    $ export OS_AUTH_URL={openstack_auth_url}
    $ export OS_USERNAME={openstack_username}
    $ export OS_PASSWORD={openstack_password}
    $ export OS_TENANT_NAME={openstack_tenant_name}


(Alternatively, download OpenStack RC file from Horizon/Access & Security/API Access.)

Source your nova credentials:

.. code-block:: console

    $ source openrc.sh


Set up your keys
----------------

Choose an existing keypair or upload a new public key, if desired.

.. code-block:: console

    $ nova keypair-add --pub-key ~/.ssh/deis.pub deis-key


Upload a CoreOS image to Glance
-------------------------------

You need to have a relatively recent CoreOS image.

If you don't already have a suitable CoreOS image and your OpenStack install allows you to upload
your own images, the following snippet will use the latest CoreOS image from the stable channel:

.. code-block:: console

    $ wget http://stable.release.core-os.net/amd64-usr/current/coreos_production_openstack_image.img.bz2
    $ bunzip2 coreos_production_openstack_image.img.bz2
    $ glance image-create --name coreos \
      --container-format bare \
      --disk-format qcow2 \
      --file coreos_production_openstack_image.img \
      --is-public True


Generate a New Discovery URL
----------------------------

.. include:: ../_includes/_generate-discovery-url.rst


Choose number of instances
--------------------------

A Deis cluster must have 3 or more nodes. See :ref:`cluster-size` for more details.

Instruct the provision script to launch the desired number of nodes:

.. code-block:: console

    $ export DEIS_NUM_INSTANCES=3


Deis network settings
---------------------

The script creates a private network called 'deis' if no such network exists.

By default, the deis subnet IP range is set to 10.21.12.0/24. To override it and the default
DNS settings, set the following variables:

.. code-block:: console

    $ export DEIS_CIDR=10.21.12.0/24
    $ export DEIS_DNS=10.21.12.3,8.8.8.8

.. note::

    This script does not handle floating IPs or routers. These should be provisioned manually by
    either Horizon or the CLI.


Run the provision script
------------------------

If you have a fairly straightforward OpenStack install, you should be able to use the provided
provisioning script. This script assumes you are using neutron and have security-groups enabled.

Run the ``provision-openstack-cluster.sh`` to spawn a new CoreOS cluster. You'll need to provide
the name of the CoreOS image name (or ID), and the key pair you just added. Optionally, you can also
specify a flavor name.

.. code-block:: console

    $ cd contrib/openstack
    $ ./provision-openstack-cluster.sh
    Usage: provision-openstack-cluster.sh <coreos image name/id> <key pair name> [flavor]
    $ ./provision-openstack-cluster.sh coreos deis-key


You can override the name of the internal network to use by setting the environment variable
``DEIS_NETWORK=internal``.  If this doesn't exist the script will try to create it with the default
CIDR which requires your OpenStack cluster to support tenant VLANs.

You can also override the name of the security group to attach to the instances by setting
``DEIS_SECGROUP=deis_test``.  If this doesn't exist the script will attempt to create it.
If you are creating your own security groups you can use the provision script as a guide.  Make sure
that you have a rule to enable full communication inside the security group, or you will have a bad day.

Manually start the instances
----------------------------

Start the instances and ensure they're operational before continuing.

Configure floating IPs
----------------------

You will want to attach a floating IP to at least one of your instances.  You'll do that like this:

.. code-block:: console

    $ nova floating-ip-create <pool>
    $ nova floating-ip-associate deis-1 <IP provided by above command>

Deploy a load balancer
----------------------

It is recommended that you deploy a load balancer for user requests to your Deis cluster.
See :ref:`configure-load-balancers` for more details on using load balancers with Deis.

Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.

Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.

.. _`contrib/openstack`: https://github.com/deis/deis/tree/master/contrib/openstack
.. _`Shlomo Hakim`: https://github.com/shakim
