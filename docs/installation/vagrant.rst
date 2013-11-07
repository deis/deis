:title: Vagrant Installation
:description: How to install a Deis controller on Vagrant for testing
:keywords: install, installation, deis, controller, setup, vagrant, virtualbox, testing

Vagrant Installation
====================

For trying out Deis, or for doing development on Deis, you can provision a
controller using Vagrant and VirtualBox. We recommend you use the binary
installer of Vagrant 1.3.5 from vagrantup.com and the 4.2.18 version of
VirtualBox.

.. include:: steps1-2.txt

3. Provision a Deis Controller
------------------------------

The ``Vagrantfile`` in the project root has the configuration for a Deis
controller. It will first need to download a base image "deis-base," which
may take a while.

Run the Vagrant provisioning script, which takes several minutes to complete:

.. code-block:: console

    $ ./contrib/ec2/provision-ec2-controller.sh us-west-2
    Creating security group: deis-controller
    + ec2-create-group deis-controller -d 'Created by Deis'
    GROUP   sg-fe82aaaa deis-controller Created by Deis
    + set +x
    Authorizing TCP ports 22,80,443,514 from 0.0.0.0/0...
    + ec2-authorize deis-controller -P tcp -p 22 -s 0.0.0.0/0
    + ec2-authorize deis-controller -P tcp -p 80 -s 0.0.0.0/0
    ...
    ec2-198-51-100-22.us-west-2.compute.amazonaws.com
    ec2-198-51-100-22.us-west-2.compute.amazonaws.com Chef Client finished, 74 resources updated
    Instance ID: i-2be2411c
    Flavor: m1.large
    Image: ami-ca63fafa
    Region: us-west-2
    Availability Zone: us-west-2b
    Security Groups: deis-controller
    Public DNS Name: ec2-198-51-100-22.us-west-2.compute.amazonaws.com
    Public IP Address: 198.51.100.22
    Run List: recipe[deis::controller]
    + set +x
    Please ensure that "deis-controller" is added to the Chef "admins" group.

.. include:: steps3-4.txt

5. Register With the Controller
-------------------------------

Registration will discover SSH keys automatically and use the
`standard environment variables`_ **AWS_ACCESS_KEY** and **AWS_SECRET_KEY** to
configure the EC2 provider with your credentials.

.. code-block:: console

    $ sudo pip install deis
    $ deis register http://deis.example.com
    username: myuser
    password:
    password (confirm):
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done

    Found EC2 credentials: AKIAJTVXXXXXXXXXXXXX
    Import these credentials? (y/n) : y
    Uploading EC2 credentials... done

6. Deploy a Formation and App
-----------------------------

Create a formation and scale it:

.. code-block:: console

    $ deis formations:create dev --flavor=ec2-us-west-2
    $ deis nodes:scale dev runtime=1

.. include:: step6.txt

.. _`Amazon EC2 API Tools`: http://aws.amazon.com/developertools/Amazon-EC2/351
.. _`standard environment variables`: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SettingUp_CommandLine.html#set_aws_credentials_linux
