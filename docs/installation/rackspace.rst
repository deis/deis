:title: Rackspace Installation
:description: How to install a Deis controller on Rackspace and create your private PaaS
:keywords: install, installation, deis, controller, setup, rackspace

Rackspace Installation
======================

.. include:: steps1-2.txt

3. Provision a Deis Controller
------------------------------

To create the controller and nodes, the Deis system expects to find a Rackspace
saved server image named "deis-base-image" in the intended cloud region. Follow
the instructions in the ``./contrib/rackspace/prepare-rackspace-image.sh``
script to create this saved image in your Rackspace account.

The `knife`_ Rackspace plugin is used to bootstrap the controller. It should
have been installed by ``bundle install``:

.. code-block:: console

    $ knife rackspace flavor list  # see if knife-rackspace works
    ID  Name                     VCPUs  RAM    Disk
    2   512MB Standard Instance  1      512    20 GB
    3   1GB Standard Instance    1      1024   40 GB
    4   2GB Standard Instance    2      2048   80 GB
    ...

Run the Rackspace provisioning script, which takes several minutes to complete:

.. code-block:: console

    $ ./contrib/rackspace/provision-rackspace-controller.sh dfw
    Creating new SSH key: deis-controller
    + ssh-keygen -f ~/.ssh/deis-controller -t rsa -N '' -C deis-controller
    + set +x
    Saved to /Users/myuser/.ssh/deis-controller
    Created data_bag[deis-users]
    Created data_bag[deis-formations]
    Created data_bag[deis-apps]
    Provisioning deis-controller with knife rackspace...
    + knife rackspace server create --bootstrap-version 11.6.2 ...
    ...
    198.51.100.22 Chef Client finished, 74 resources updated
    198.51.100.22
    + set +x
    Please ensure that "deis-controller" is added to the Chef "admins" group.

.. include:: steps3-4.txt

5. Register With the Controller
-------------------------------

Registration will discover SSH keys automatically and use the environment
variables **RACKSPACE_USERNAME** and **RACKSPACE_API_KEY** to configure
the Rackspace provider with your credentials.

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

    Found Rackspace credentials: 8b7e2cXXXXXXXXXXXXXXXXXXXXXXXXX
    Import these credentials? (y/n) : y
    Uploading Rackspace credentials... done

6. Deploy a Formation and App
-----------------------------

Create a formation and scale it:

.. code-block:: console

    $ deis formations:create dev --flavor=rackspace-dfw
    $ deis nodes:scale dev runtime=1

.. include:: step6.txt
