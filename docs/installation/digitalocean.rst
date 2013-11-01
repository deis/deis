:title: DigitalOcean Installation
:description: How to install a Deis controller on DigitalOcean and create your private PaaS
:keywords: install, installation, deis, controller, setup, digitalocean

DigitalOcean Installation
=========================

.. include:: steps1-2.txt

3. Provision a Deis Controller
------------------------------

To create the controller and nodes, the Deis system expects to find a
DigitalOcean snapshot (or droplet image) named "deis-base" in the intended
cloud region.

The `knife`_ DigitalOcean plugin is used to bootstrap the controller. It should
have been installed by ``bundle install``:

.. code-block:: console

    $ knife digital_ocean region list  # list regions, test knife-DO
    ID  Name
    1   New York 1
    2   Amsterdam 1
    3   San Francisco 1
    4   New York 2

The DigitalOcean provisioning script expects one argument: the ID of the cloud
region in which to host the controller. Take note of the output of the
test command above, then follow the instructions in the
``./contrib/digitalocean/prepare-digitalocean-snapshot.sh`` script to create a
Deis-compatible snapshot in your DigitalOcean account.

Run the DigitalOcean provisioning script, which takes several minutes to complete:

.. code-block:: console

    $ ./contrib/digitalocean/provision-digitalocean-controller.sh 4
    Provisioning a deis controller on Digital Ocean!
    Creating new SSH key: deis-controller
    + ssh-keygen -f ~/.ssh/deis-controller -t rsa -N '' -C deis-controller
    ...
    Created data_bag[deis-apps]
    Provisioning deis-controller with knife digital_ocean...
    + knife digital_ocean droplet create --bootstrap-version 11.6.2 ...
    Droplet creation for deis-controller started. Droplet-ID is 123456
    Waiting for IPv4-Address.done
    IPv4 address is: 198.51.100.22
    ...
    198.51.100.22 Chef Client finished, 74 resources updated
    198.51.100.22
    + set +x
    Please ensure that "deis-controller" is added to the Chef "admins" group.

.. include:: steps3-4.txt

5. Register With the Controller
-------------------------------

Registration will discover SSH keys automatically and use the environment
variables **DIGITALOCEAN_CLIENT_ID** and **DIGITALOCEAN_API_KEY** to configure
the DigitalOcean provider with your credentials.

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

    Found Digitalocean credentials: hkrVAMXXXXXXXXXXXXXXXX
    Import these credentials? (y/n) : y
    Uploading Digitalocean credentials... done

6. Deploy a Formation and App
-----------------------------

Create a formation and scale it:

.. code-block:: console

    $ deis formations:create dev --flavor=digitalocean-new-york-2
    $ deis nodes:scale dev runtime=1

.. include:: step6.txt
