:title: Installing Deis on AWS
:description: How to provision a multi-node Deis cluster on Amazon AWS

.. _deis_on_aws:

Amazon AWS
==========

In this tutorial, we will show you how to set up your own 3-node cluster on Amazon Web Services.

Please :ref:`get the source <get_the_source>` and refer to the scripts in `contrib/ec2`_
while following this documentation.


Install the AWS Command Line Interface
--------------------------------------

In order to start working with Amazon's API, let's install `awscli`_:

.. code-block:: console

    $ pip install awscli

We'll also need `PyYAML`_ for the Deis EC2 provision script to run:

.. code-block:: console

    $ pip install pyyaml


Configure aws-cli
-----------------

Run ``aws configure`` to set your AWS credentials:


.. code-block:: console

    $ aws configure
    AWS Access Key ID [None]: ***************
    AWS Secret Access Key [None]: ************************
    Default region name [None]: us-west-1
    Default output format [None]:


Upload keys
-----------

Generate and upload a new keypair to AWS, ensuring that the name of the keypair is set to "deis".

.. code-block:: console

    $ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis
    $ aws ec2 import-key-pair --key-name deis --public-key-material file://~/.ssh/deis.pub


Choose Number of Instances
--------------------------

By default, the script will provision 3 servers. You can override this by setting
``DEIS_NUM_INSTANCES``:

.. code-block:: console

    $ export DEIS_NUM_INSTANCES=5

Note that for scheduling to work properly, clusters must consist of at least 3 nodes and always
have an odd number of members. For more information, see `etcd disaster recovery`_.

Deis clusters of less than 3 nodes are unsupported.


Generate a New Discovery URL
----------------------------

.. include:: ../_includes/_generate-discovery-url.rst


Customize cloudformation.json
-----------------------------

Any of the parameter defaults defined in deis.template.json can be overridden by setting the value
in `cloudformation.json`_ like so:

.. code-block:: console

    {
        "ParameterKey":     "InstanceType",
        "ParameterValue":   "m3.xlarge"
    },
    {
        "ParameterKey":     "KeyPair",
        "ParameterValue":   "jsmith"
    },
    {
        "ParameterKey":     "EC2VirtualizationType",
        "ParameterValue":   "PV"
    },
    {
        "ParameterKey":     "AssociatePublicIP",
        "ParameterValue":   "false"
    }

The only entry in cloudformation.json required to launch your cluster is `KeyPair`, which is
already filled out. The defaults will be applied for the other settings.

If updated with update-ec2-cluster.sh, the InstanceType will only impact newly deployed instances
(`#1758`_).

NOTE: The smallest recommended instance size is `large`. Having not enough CPU or RAM will result
in numerous issues when using the cluster.


Launch into an existing VPC
---------------------------

By default, the provided CloudFormation script will create a new VPC for Deis. However, the script
supports provisioning into an existing VPC instead. You'll need to have a VPC configured with an
internet gateway and a sane routing table (the default VPC in a region should be ready to go).

To launch your cluster into an existing VPC, export three additional environment variables:

 - ``VPC_ID``
 - ``VPC_SUBNETS``
 - ``VPC_ZONES``

``VPC_ZONES`` must list the availability zones of the subnets in order.

For example, if your VPC has ID ``vpc-a26218bf`` and consists of the subnets ``subnet-04d7f942``
(which is in ``us-east-1b``) and ``subnet-2b03ab7f`` (which is in ``us-east-1c``) you would export:

.. code-block:: console

    export VPC_ID=vpc-a26218bf
    export VPC_SUBNETS=subnet-04d7f942,subnet-2b03ab7f
    export VPC_ZONES=us-east-1b,us-east-1c


Run the Provision Script
------------------------

Run the cloudformation provision script to spawn a new CoreOS cluster:

.. code-block:: console

    $ cd contrib/ec2
    $ ./provision-ec2-cluster.sh
    {
        "StackId": "arn:aws:cloudformation:us-west-1:413516094235:stack/deis/9699ec20-c257-11e3-99eb-50fa01cd4496"
    }
    Your Deis cluster has successfully deployed.
    Please wait for all instances to come up as "running" before continuing.

.. note::

    The default name of the CloudFormation stack will be ``deis``. You can specify a different name
    with ``./provision-ec2-cluster.sh <name>``.

Check the AWS EC2 web control panel and wait until "Status Checks" for all instances have passed.
This will take several minutes.


Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.

CloudFormation Updates
----------------------

To use CloudFormation to perform update operations to your stack, there is another script: 
`update_ec2_cluster.sh`_. Depending on the parameters that you have changed, CloudFormation
may replace the EC2 instances in your stack.

The following parameters can be changed without replacing all instances in a stack:

- ``ClusterSize`` - Number of nodes in the cluster. This may launch new instances or terminate
  existing instances. If you are scaling down, this may interrupt service. If a container
  was running on an instance that was terminated, it will have to be rebalanced onto another 
  node which will cause some downtime.
- ``SSHFrom`` - Locks down SSH access to the Deis hosts. This will update the security
  group for the Deis hosts.

Please reference the AWS documentation for `more information about CloudFormation stack updates`_.

.. _`#1758`: https://github.com/deis/deis/issues/1758
.. _`awscli`: https://github.com/aws/aws-cli
.. _`contrib/ec2`: https://github.com/deis/deis/tree/master/contrib/ec2
.. _`cloudformation.json`: https://github.com/deis/deis/blob/master/contrib/ec2/cloudformation.json
.. _`etcd`: https://github.com/coreos/etcd
.. _`etcd disaster recovery`: https://github.com/coreos/etcd/blob/master/Documentation/admin_guide.md#disaster-recovery
.. _`PyYAML`: http://pyyaml.org/
.. _`update_ec2_cluster.sh`: https://github.com/deis/deis/blob/master/contrib/ec2/update-ec2-cluster.sh
.. _`More information about CloudFormation stack updates`: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks.html

