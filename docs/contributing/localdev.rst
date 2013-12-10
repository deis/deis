:description: Local development setup instructions for contributing to the Deis project.
:keywords: deis, documentation, contributing, developer, setup, chef, knife, vagrant

.. _localdev:

Local Development
=================
We have tried to make it simple to hack on Deis, but as an open PaaS 
there are necessarily several moving pieces and some setup required. 
We welcome any suggestions for automating or simplifying this process.

Prerequisites
-------------
We strongly recommend using `Vagrant`_ with `VirtualBox`_ so you can 
develop inside a set of isolated virtual machines. You will need:

 * Vagrant 1.3.5+
 * VirtualBox 4.2+

Fork the Repository
-------------------
To get Deis running locally, first `fork the Deis repository`_ at GitHub.com.
Then clone your fork of the repository for local development:

.. code-block:: console

	$ git clone https://github.com/<username>/deis.git
	$ cd deis
	$ export DEIS_DIR=`pwd`  # to use in future commands

Configure a Chef Server
-----------------------
Deis relies on Chef Server to provide a battle-tested foundation for
hosting `cookbooks`_ used to configure :ref:`Nodes <node>` and
`data bags`_ used to store cluster configuration.
You can run a local Chef Server or use a Hosted Chef Server.

Local Chef Server
`````````````````
For development purposes, you can spin up a local Chef Server using Vagrant.
Please note it will take up at least 1GB of RAM.  From your workstation:

.. code-block:: console

    $ cd $DEIS_DIR/contrib/vagrant/chef-server
    $ vagrant plugin install vagrant-vbguest
    $ ln -s $DEIS_DIR/contrib/vagrant/knife-config $DEIS_DIR/.chef
    $ vagrant up
    $ vagrant ssh -c 'sudo cat /etc/chef-server/admin.pem' > $DEIS_DIR/.chef/admin.pem
    $ vagrant ssh -c 'sudo cat /etc/chef-server/chef-validator.pem' > $DEIS_DIR/.chef/chef-validator.pem
    $ echo '{"ssl":{"verify": false }}' > ~/.berkshelf/config.json

Now use ``knife client list`` to test connectivity to the local Chef Server:

.. code-block:: console

    $ knife client list
    chef-validator
    chef-webui

Hosted Chef Server
``````````````````
If you don't want to run your own Chef server, you can
`sign up for a free Hosted Chef account`_.

 * `Login to the Chef Server <https://preview.opscode.com/login>`_
 * Click on the ``Administration`` tab and choose your organization
 * Click ``Starter Kit`` in the sidebar to start a download
 * Copy the ``.chef`` directory inside the Starter Kit into the root of your Deis checkout

Now use ``knife client list`` to test connectivity to the local Chef Server:

.. code-block:: console

    $ knife client list
    gabrtv-validator

You should see at least the validator key for your Chef organization.
If not, your `knife.rb`_ configuration or Chef keys are probably incorrect.

Upload Cookbooks
----------------
Upload the current Deis cookbooks using Berkshelf:

.. code-block:: console

    $ gem install bundler  # install the bundler tool (if necessary)
    $ bundle install       # install ruby dependencies from Gemfile
    $ berks install        # install cookbooks to your local berkshelf
    $ berks upload         # upload berkshelf cookbooks to the chef server

Provision the Controller
------------------------
Now that the Chef Server is in place with the latest version of our cookbooks,
we can provision the :ref:`controller`.

.. code-block:: console

    $ cd $DEIS_DIR/contrib/vagrant
    $ ./provision-controller.sh

The provisioning process will ask a few questions and then run a
``knife bootstrap`` of the controller.

Add Controller to Admin Group
`````````````````````````````
In order for the controller to delete records on the Chef Server,
it must be part of the Admin group.

For a Local Chef Server

 * Open a shell in the project repository
 * Run ``knife client edit deis-controller``
 * Set "admin" to "true" and save the file

For a Hosted Chef Server

 * Login to the Web UI
 * Navigate to the Groups tab 
 * Click "Edit" on the "Admins" row
 * Under "Clients" heading, toggle the "deis-controller" radio button
 * Save changes

SSH Access
``````````
The Controller needs to be able to run Vagrant commands on your host machine. 
It does this via SSH. Therefore you will need a running SSH server open on port 22
and a means to broadcast your hostname to local DNS.

 * On Mac OSX you just need to go to "System Preferences > Sharing" and enable "Remote Login"
 * On Debian-flavoured Linux you just need to ``sudo apt-get install openssh-server avahi-daemon``

Install the Client
------------------
In a development environment you'll want to use the latest version of the client.
Install it by symlinking ``client/deis.py`` to ``deis`` on your local workstation.

.. code-block:: console

    $ ln -fs $DEIS_DIR/client/deis.py /usr/local/bin/deis
    $ deis
    Usage: deis <command> [<args>...]

Register an Admin User
----------------------
Use the Deis client to register a new user on the controller.
As the first user, you will receive full admin permissions.

.. code-block:: console

    $ deis register http://deis-controller.local
    username: myuser
    password:
    password (confirm):
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

Once the user is registered, activate the Vagrant provider with:

.. code-block:: console

    $ deis providers:discover
    No EC2 credentials discovered.
    No Rackspace credentials discovered.
    No DigitalOcean credentials discovered.
    Discovered locally running Deis Controller VM
    Activating Vagrant as a provider... done

Add your SSH key for ``git push`` access using:

.. code-block:: console

    $ deis keys:add
    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done

Deploy a Vagrant Formation
--------------------------
These are 3 default flavors of Vagrant nodes: 512MB, 1024MB and 2048MB.
To create a formation with a 512MB nodes:

.. code-block:: console

    $ deis formations:create dev --flavor=vagrant-512
    $ deis nodes:scale dev runtime=1

This will use the Deis :ref:`Provider` API to spin up a new Vagrant node as
part of a single-host formation.  The scaling process can take ~ 5 min
as Vagrant boots a host and runs through the first Chef converge.

Once ``nodes:scale`` returns, your local development environment is running!
Follow the rest of the :ref:`Developer Guide <developer>` to 
deploy your first application.

Useful Commands
---------------
Once your controller is running, here are some helpful commands.

Tail Logs
`````````
.. code-block:: console

    $ vagrant ssh -c 'sudo tail -f /var/log/upstart/deis-* /var/log/deis/*'

Restart Services
````````````````
.. code-block:: console

    $ vagrant ssh -c 'sudo restart deis-worker && sudo restart deis-server'

Django Admin
````````````
.. code-block:: console

    $ vagrant ssh              # SSH into the controller
    $ sudo su deis -l          # change to deis user
    $ cd controller            # change into the django project root
    $ source venv/bin/activate # activate python virtualenv
    $ ./manage.py shell        # get a django shell

Have commands other Deis developers might find useful? Send us a PR!

Standards & Test Coverage
-------------------------
When changing Python code in the Deis project, keep in mind our :ref:`standards`.
Specifically, when you change local code, you must run
``make flake8 && make coverage``, then check the HTML report to see
that test coverage has improved as a result of your changes and new unit tests.

.. code-block:: console

	$ make flake8
	flake8
	./api/models.py:17:1: F401 'Group' imported but unused
	./api/models.py:81:1: F841 local variable 'result' is assigned to but never used
	make: *** [flake8] Error 1
	$
	$ make coverage
	coverage run manage.py test api celerytasks client web
	Creating test database for alias 'default'...
	...................ss
	----------------------------------------------------------------------
	Ran 21 tests in 18.135s

	OK (skipped=2)
	Destroying test database for alias 'default'...
	coverage html
	$ head -n 25 htmlcov/index.html | grep pc_cov
	            <span class='pc_cov'>81%</span>

Pull Requests
-------------
Please create a GitHub `pull request`_ for any code changes that will benefit Deis users
in general. This workflow helps changesets map well to discrete features.

Creating a pull request on the Deis repository also runs a Travis CI build to
ensure the pull request doesn't break any tests or reduce code coverage.

Cookbook Development
--------------------
If you want to modify Deis' Chef recipes, you should also clone the `deis-cookbook`_
repository:

.. code-block:: console

	$ git clone -q https://github.com/opdemand/deis-cookbook.git

Please see `deis-cookbook`_ for information about contributing Chef code to Deis.

.. _`Vagrant`: http://www.vagrantup.com/
.. _`VirtualBox`: https://www.virtualbox.org/
.. _`sign up for a free Hosted Chef account`: https://getchef.opscode.com/signup
.. _`knife.rb`: http://docs.opscode.com/config_rb_knife.html
.. _`cookbooks`: http://docs.opscode.com/essentials_cookbooks.html
.. _`data bags`: http://docs.opscode.com/essentials_data_bags.html
.. _`fork the Deis repository`: https://github.com/opdemand/deis/fork
.. _`deis-cookbook`: https://github.com/opdemand/deis-cookbook
.. _`pull request`: https://github.com/opdemand/deis/pulls
