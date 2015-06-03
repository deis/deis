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


During installation, ``deisctl`` will make an SSH connection to the cluster.
It will need to be able to use this key to connect.

Most users use SSH agent (``ssh-agent``). If this is the case, run
``ssh-agent add ~/.ssh/deis`` to add the key. Otherwise, you may prefer to
modify ``~/.ssh/config`` to add the key to the IPs in AWS.

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

The configuration files and templates for AWS are located in the directory
``contrib/ec2/`` in the Deis repository.

Any of the parameter defaults defined in ``deis.template.json`` can be
overridden by setting the value in `cloudformation.json`_. For example, to
configure all of the options to non-default values:

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
    },
    {
        "ParameterKey":     "ELBScheme",
        "ParameterValue":   "internal"
    },
    {
        "ParameterKey":     "RootVolumeSize",
        "ParameterValue":   "100"
    },
    {
        "ParameterKey":     "DockerVolumeSize",
        "ParameterValue":   "1000"
    },
    {
        "ParameterKey":     "EtcdVolumeSize",
        "ParameterValue":   "5"
    }


The only entry in cloudformation.json required to launch your cluster is `KeyPair`, which is
already filled out. The defaults will be applied for the other settings. The default values are
defined in ``deis.template.json``.

If updated with ``update-ec2-cluster.sh``, the InstanceType will only impact newly deployed instances
(`#1758`_).

NOTE: The smallest recommended instance size is ``large``. Having not enough CPU or RAM will result
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

If you have set up private subnets in which you'd like to run your Deis hosts, and public subnets
for the ELB, you should export the following environment variables instead:

 - ``VPC_ID``
 - ``VPC_SUBNETS``
 - ``VPC_PRIVATE_SUBNETS``
 - ``VPC_ZONES``

For example, if you have a public subnet ``subnet-8cd457b3`` for the ELB and a private subnet
``subnet-8cd457b0`` (both in ``us-east-1a``) you would export:

.. code-block:: console

    export VPC_ID=vpc-a26218bf
    export VPC_SUBNETS=subnet-8cd457b3
    export VPC_PRIVATE_SUBNETS=subnet-8cd457b0
    export VPC_ZONES=us-east-1a


Run the Provision Script
------------------------

Run the cloudformation provision script to spawn a new CoreOS cluster:

.. code-block:: console

    $ cd contrib/ec2
    $ ./provision-ec2-cluster.sh
    Creating CloudFormation stack deis
    {
        "StackId": "arn:aws:cloudformation:us-east-1:69326027886:stack/deis/1e9916b0-d7ea-11e4-a0be-50d2020578e0"
    }
    Waiting for instances to be created...
    Waiting for instances to be created... CREATE_IN_PROGRESS
    Waiting for instances to pass initial health checks...
    Waiting for instances to pass initial health checks...
    Waiting for instances to pass initial health checks...
    Instances are available:
    i-5c3c91aa	203.0.113.91	m3.large	us-east-1a	running
    i-403c91b6	203.0.113.20	m3.large	us-east-1a	running
    i-e36fc6ee	203.0.113.31	m3.large	us-east-1b	running
    Using ELB deis-DeisWebE-17PGCR3KPJC54 at deis-DeisWebE-17PGCR3KPJC54-1499385382.us-east-1.elb.amazonaws.com
    Your Deis cluster has been successfully deployed to AWS CloudFormation and is started.
    Please continue to follow the instructions in the documentation.

.. note::

    The default name of the CloudFormation stack will be ``deis``. You can specify a different name
    with ``./provision-ec2-cluster.sh <name>``.

Remote IPs behind your ELB
--------------------------

The ELB you just created is load-balancing raw TCP connections, which is required for custom domain SSL
and WebSockets. As remote IPs are by default not visible behind a TCP-Proxy, the ELB and your cluster routers
were created with `Proxy Protocol`_ enabled.


Configure DNS
-------------

You will need a DNS entry that points to the ELB instance created above. Find
the ELB name in the AWS web console or by running ``aws elb describe-load-balancers``
and finding the Deis ELB.

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
.. _`Proxy Protocol`: http://docs.aws.amazon.com/ElasticLoadBalancing/latest/DeveloperGuide/enable-proxy-protocol.html
