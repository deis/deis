:title: Flavor
:description: A Deis flavor defines the configuration for nodes in a layer, including their provider type and launch parameters.
:keywords: flavor, deis, nodes, configuration

.. _flavor:

Flavor
======
A flavor defines the configuration for :ref:`Nodes <node>` in a :ref:`Layer`.
Flavors are used by the :ref:`Provider` API to customize cloud server configuration.
Configuration parameters include region/location, server size and OS image.

Default Flavors
---------------
Each Deis user account contains a set of default :ref:`Flavors <flavor>` 
for each :ref:`provider`.  Typically, there is one flavor for each region/location.
In the case of the EC2 provider, the default flavors include:

 * ec2-us-east-1
 * ec2-us-west-1
 * ec2-us-west-2
 * ec2-eu-west-1
 * ec2-ap-northeast-1
 * ec2-ap-southeast-1
 * ec2-ap-southeast-2
 * ec2-sa-east-1

Each of these default flavors uses a specific instance size and an optimized EC2 AMI
that speeds boot times during scale operations.  Other providers come with default
flavors that are similarly optimized.

Creating a Flavor
-----------------
Let's create a new flavor from scratch using the EC2 provider.  Our goals are:

 #. Use the us-east-1 region
 #. Use a custom base image (EC2 AMI)
 #. Use the m1.large instance type

To create the flavor we'll use the "deis flavors:create" command.  
We must pass it a name for the flavor (ec2-custom), the provider type (ec2)
and some JSON that defines the configuration we want.
 
.. code-block:: console

    $ deis flavors:create ec2-custom --provider=ec2 --params='
    > {"region": "us-east-1", "image": "ami-fa99c193", "size": "m1.large"}'
    ec2-custom

.. info:
   Each provider supports different JSON fields.  Consult the provider
   documentation under Server Reference for more details.

Nice work!  We now have an "ec2-custom" flavor built to our specifications.

Viewing a Flavor
----------------
To make sure a flavor is correct, we can use "deis flavors:info".

.. code-block:: console

    $ deis flavors:info ec2-custom
    {
      "updated": "2013-11-26T22:15:34.520Z", 
      "uuid": "86213a3e-62f6-4d9b-91a9-7cf65c33ba6f", 
      "created": "2013-11-26T22:15:34.520Z", 
      "params": "{\"region\": \"us-east-1\", \"image\": \"ami-fa99c193\", \"size\": \"m1.large\"}", 
      "provider": "ec2", 
      "owner": "gabrtv", 
      "id": "ec2-custom"
    }

Updating a Flavor
-----------------
Let's go back and hardcode a zone in our "ec2-custom" flavor.  The process is similar:

.. code-block:: console

    $ deis flavors:update ec2-custom '{"zone": "us-east-1a"}'
    {
      "updated": "2013-11-26T22:25:17.260Z", 
      "uuid": "86213a3e-62f6-4d9b-91a9-7cf65c33ba6f", 
      "created": "2013-11-26T22:15:34.520Z", 
      "params": "{\"region\": \"us-east-1\", \"image\": \"ami-fa99c193\", \"zone\": \"us-east-1a\", \"size\": \"m1.large\"}", 
      "provider": "ec2", 
      "owner": "gabrtv", 
      "id": "ec2-custom"
    }

Great.  We now have our zone inside the "params" field.
Notice how the other fields were kept even though we didn't specify, for example, a region.
This behavior makes it easy to update flavors in place.

Deleting a Flavor
-----------------
Before you delete a flavor, make sure servers in Deis-land aren't using it.
The "deis flavors:delete" command is simple.

.. code-block:: console

    $ deis flavors:delete ec2-custom
