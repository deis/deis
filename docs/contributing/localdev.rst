:description: Local development setup instructions for contributing to the Deis project.

.. _localdev:

Local Development
=================

We try to make it simple to hack on Deis. However, there are necessarily several moving
pieces and some setup required. We welcome any suggestions for automating or simplifying
this process.

If you're just getting into the Deis codebase, look for GitHub issues with the label
`easy-fix`_. These are more straightforward or low-risk issues for which we need pull
requests. Issues tagged `easy-fix`_ are a great way to become more familiar with Deis.

Prerequisites
-------------

We strongly recommend using `Vagrant`_ with `VirtualBox`_ so you can develop inside a set
of isolated virtual machines. You will need:

 * Vagrant 1.6+
 * VirtualBox 4.2+

Fork the Repository
-------------------

To get Deis running locally, first `fork the Deis repository`_, then clone your fork of
the repository for local development:

.. code-block:: console

	$ git clone git@github.com:<username>/deis.git
	$ cd deis
	$ export DEIS_DIR=`pwd`  # to use in future commands

Provision the Cluster
---------------------

First bring up a virtual machine to host Deis. To share your local codebase into the
CoreOS VM, Deis uses rsync.

.. code-block:: console

    $ vagrant up
    Bringing machine 'deis-1' up with 'virtualbox' provider...
    ==> deis-1: Importing base box 'coreos-402.2.0'...
    ==> deis-1: Matching MAC address for NAT networking...
    ==> deis-1: Setting the name of the VM: deis_deis-1_1408039741706_8568
    ==> deis-1: Clearing any previously set network interfaces...
    ==> deis-1: Preparing network interfaces based on configuration...
        deis-1: Adapter 1: nat
        deis-1: Adapter 2: hostonly
    ==> deis-1: Forwarding ports...
        deis-1: 22 => 2222 (adapter 1)
    ==> deis-1: Running 'pre-boot' VM customizations...
    ==> deis-1: Booting VM...
    ==> deis-1: Waiting for machine to boot. This may take a few minutes...
    $

Set ``FLEETCTL_TUNNEL`` so the ``fleetctl`` client on your workstation can connect to the VM:

.. code-block:: console

    export FLEETCTL_TUNNEL=172.17.8.100

Next, run ``make pull && make build`` to SSH into the VM, pull Deis' images from the
Docker Index, then update those images with any local changes.

.. code-block:: console

    $ make pull
    for host in 172.17.8.100 ; do ssh -o LogLevel=FATAL -o Compression=yes \
        -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -o PasswordAuthentication=no core@$host \
        -t 'for c in builder cache controller database logger registry router; do docker pull deis/$c:latest; done'; done
    Pulling repository deis/registry
    d2c347aa26dd: Pulling dependent layers
    511136ea3c5a: Download complete
    6170bb7b0ad1: Download complete
    79fdb1362c84: Downloading [====>
    ...
    e5efa1477310: Download complete
    Connection to 127.0.0.1 closed.
    $ make build
    for host in 172.17.8.100 ; do ssh -o LogLevel=FATAL -o Compression=yes \
        -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
        -o PasswordAuthentication=no core@$host \
        -t 'for c in builder cache controller database logger registry router; do cd $c && docker build -t deis/$c . && cd ..; done'; done
    Uploading context 22.53 kB
    Uploading context
    Step 0 : FROM deis/base:latest
    Pulling repository deis/base
    60024338bc63: Download complete
    ...
    Step 12 : CMD ["/app/bin/boot"]
     ---> Running in ccdc3d283f4f
     ---> cf4b7a398500
    Removing intermediate container ccdc3d283f4f
    Successfully built cf4b7a398500
    Connection to 127.0.0.1 closed.

Finally, do ``make run`` to start all Deis containers:

.. code-block:: console

    $ make run
    Job deis-router@1.service loaded on ff3442c2.../172.17.8.100
    Job deis-builder-data.service loaded on ff3442c2.../172.17.8.100
    Job deis-database-data.service loaded on ff3442c2.../172.17.8.100
    Job deis-logger-data.service loaded on ff3442c2.../172.17.8.100
    Job deis-registry-data.service loaded on ff3442c2.../172.17.8.100
    fleetctl --strict-host-key-checking=false load logger/systemd/deis-logger.service cache/systemd/deis-cache.service database/systemd/deis-database.service
    Job deis-cache.service loaded on ff3442c2.../172.17.8.100
    Job deis-database.service loaded on ff3442c2.../172.17.8.100
    Job deis-logger.service loaded on ff3442c2.../172.17.8.100
    fleetctl --strict-host-key-checking=false load registry/systemd/*.service
    Job deis-registry.service loaded on ff3442c2.../172.17.8.100
    fleetctl --strict-host-key-checking=false load controller/systemd/*.service
    Job deis-controller.service loaded on ff3442c2.../172.17.8.100
    fleetctl --strict-host-key-checking=false load builder/systemd/*.service
    Job deis-builder.service loaded on ff3442c2.../172.17.8.100
    Deis components may take a long time to start the first time they are initialized.
    Waiting for 1 of 1 deis-routers to start...

Install the Client
------------------

In a development environment you'll want to use the latest version of the client. Install
its dependencies by using the Makefile and symlinking ``client/deis.py`` to ``deis`` on
your local workstation.

.. code-block:: console

    $ cd $DEIS_DIR/client
    $ make install
    $ ln -fs $DEIS_DIR/client/deis.py /usr/local/bin/deis
    $ deis
    Usage: deis <command> [<args>...]

Register an Admin User
----------------------

Use the Deis client to register a new user on the controller. As the first user, you will
receive full admin permissions.

.. code-block:: console

    $ deis register http://deis.local.deisapp.com
    username: myuser
    password:
    password (confirm):
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

Once the user is registered, add your SSH key for ``git push`` access using:

.. code-block:: console

    $ deis keys:add
    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /home/myuser/.ssh/id_rsa.pub to Deis... done

Your local development environment is running! Follow the rest of the `_using_deis` guide
to deploy your first application.

Test Your Changes
-----------------

In the single-node Vagrant environment, testing your changes to Deis itself is easy!

- Make changes to the code in one of the component subdirectories, such as
  ``controller/``
- run ``make -C controller/ build run``
- Test your changes with ``make -C controller/ test-unit`` and interactively with the
  Deis client

Useful Commands
---------------

Once your controller is running, here are some helpful commands.

Tail Logs
`````````

.. code-block:: console

    $ vagrant ssh -c 'docker logs -f deis-controller'

Rebuild Services from Source
````````````````````````````

    $ make -C controller build

Restart Services
````````````````

.. code-block:: console

    $ make -C controller restart

Django Admin
````````````

.. code-block:: console

    $ vagrant ssh              # SSH into the controller
    $ nse deis-controller      # inject yourself into the container
    $ cd /app                  # change into the django project root
    $ ./manage.py shell        # get a django shell

Have commands other Deis developers might find useful? Send us a PR!

Standards & Test Coverage
-------------------------

When changing Python code in the Deis project, keep in mind our :ref:`standards`.
Specifically, when you change local code, you must run ``make flake8 && make coverage``,
then check the HTML report to see that test coverage has improved as a result of your
changes and new unit tests.

.. code-block:: console

	$ make flake8
	flake8
	./api/models.py:17:1: F401 'Group' imported but unused
	./api/models.py:81:1: F841 local variable 'result' is assigned to but never used
	make: *** [flake8] Error 1
	$
	$ make coverage
	coverage run manage.py test --noinput api web
	WARNING Cannot synchronize with etcd cluster
	Creating test database for alias 'default'...
	...............................................
	----------------------------------------------------------------------
	Ran 47 tests in 47.768s

	OK
	Destroying test database for alias 'default'...
	coverage html
	$ head -n 25 htmlcov/index.html | grep pc_cov
	            <span class='pc_cov'>81%</span>

Pull Requests
-------------

Please create a GitHub `pull request`_ for any code changes that will benefit Deis users
in general. This workflow helps changesets map well to discrete features.

Creating a pull request on the Deis repository also runs an integration test on
http://ci.deis.io to ensure the pull request doesn't break any tests or reduce code
coverage.


.. _`easy-fix`: https://github.com/deis/deis/issues?labels=easy-fix&state=open
.. _`Vagrant`: http://www.vagrantup.com/
.. _`VirtualBox`: https://www.virtualbox.org/
.. _`fork the Deis repository`: https://github.com/deis/deis/fork
.. _`pull request`: https://github.com/deis/deis/pulls
