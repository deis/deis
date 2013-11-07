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
controller. Vagrant will first need to download a base image "deis-base," which
may take a while.

Run the Vagrant provisioning script, which takes several minutes to complete.
Pay attention, because it will ask for confirmation that an SSH key can be
added to your $HOME/.ssh/authorized_keys file:

.. code-block:: console

    $ ./contrib/vagrant/provision-controller.sh
    Created data_bag[deis-users]
    Created data_bag[deis-formations]
    Created data_bag[deis-apps]
    Booting deis-controller with 'vagrant up'
    ~/Projects/deis ~/Projects/deis
    Bringing machine 'default' up with 'virtualbox' provider...
    [default] Importing base box 'deis-node'...
    [default] Matching MAC address for NAT networking...
    ...
    [default] Running: inline script
    stdin: is not a tty
    avahi-daemon stop/waiting
    avahi-daemon start/running, process 1366
    Add the Deis Controller's SSH key to your authorized_keys file? y
    Generating public/private rsa key pair.
    Your identification has been saved in /home/vagrant/.ssh/id_rsa.
    ...
    deis-controller.local     - execute "bash"  "/tmp/chef-script20131107-1476-la5wbp"
    deis-controller.local
    deis-controller.local
    deis-controller.local
    deis-controller.local Chef Client finished, 77 resources updated
    deis-controller.local
    + set +x
    Updating Django site object from 'example.com' to 'deis-controller'...
    Site object updated.
    ~/Projects/deis
    Please ensure that "deis-controller" is added to the Chef "admins" group.

.. include:: steps3-4.txt

5. Register With the Controller
-------------------------------

Registration will discover the local Deis controller running in Vagrant and
set up the necessary provider entry so that the controller can SSH back to
the host, which is necessary to run "vagrant up" and thus scale nodes.

.. code-block:: console

    $ sudo pip install deis
    $ deis register http://deis-controller.local
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

    Detected locally running Deis Controller VM
    Activating Vagrant as a provider... done

6. Deploy a Formation and App
-----------------------------

Create a formation and scale it:

.. code-block:: console

    $ deis formations:create dev --flavor=vagrant-1024
    $ deis nodes:scale dev runtime=1

.. include:: step6.txt

