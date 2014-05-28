:description: Local development setup instructions for contributing to the Deis project.

.. _localdev:

Local Development
=================
We try to make it simple to hack on Deis, but as an open PaaS, there are
necessarily several moving pieces and some setup required. We welcome
any suggestions for automating or simplifying this process.

If you're just getting into the Deis codebase, look for GitHub issues
with the label `easy-fix`_. These are more straightforward or low-risk
issues for which we need pull requests. Issues tagged `easy-fix`_ are a
great way to become more familiar with Deis.

Prerequisites
-------------
We strongly recommend using `Vagrant`_ with `VirtualBox`_ so you can
develop inside a set of isolated virtual machines. You will need:

 * Vagrant 1.5.1+
 * VirtualBox 4.2+

Fork the Repository
-------------------
To get Deis running locally, first `fork the Deis repository`_, then
clone your fork of the repository for local development:

.. code-block:: console

	$ git clone git@github.com:<username>/deis.git
	$ cd deis
	$ export DEIS_DIR=`pwd`  # to use in future commands

Provision the Controller
------------------------
First bring up a virtual machine to host Deis. To share your local
codebase into the CoreOS VM, Deis uses NFS mounts, so you will be
prompted for an administrative password.

.. code-block:: console

    $ vagrant up
    Bringing machine 'deis' up with 'virtualbox' provider...
    ==> deis: Importing base box 'coreos-alpha'...
    ...
    ==> deis: Exporting NFS shared folders...
    ==> deis: Preparing to edit /etc/exports. Administrator privileges will be required...
    Password:
    ==> deis: Mounting NFS shared folders...
    ...
    ==> deis: Running provisioner: shell...
    deis: Running: inline script
    $

Set environment variables to have the `docker` and `fleetctl` clients on
your workstation connect to the VM:

.. code-block:: console

    export DOCKER_HOST=tcp://172.17.8.100:4243
    export FLEETCTL_TUNNEL=172.17.8.100

Next, run ``make pull && make build`` to SSH into the VM, pull Deis'
images from the Docker Index, then update those images with any local
changes.

.. code-block:: console

    $ make pull
    vagrant ssh -c 'for c in registry logger database cache controller \
      builder router; do docker pull deis/$c; done'
    Pulling repository deis/registry
    d2c347aa26dd: Pulling dependent layers
    511136ea3c5a: Download complete
    6170bb7b0ad1: Download complete
    79fdb1362c84: Downloading [====>
    ...
    e5efa1477310: Download complete
    Connection to 127.0.0.1 closed.
    $ make build
    vagrant ssh -c 'cd share && for c in registry logger database \
      cache controller builder router; \
      do cd $c && docker build -t deis/$c . && cd ..; done'
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

Finally, do ``make run`` to start all Deis containers and displays their
collected log output:

.. code-block:: console

    $ make run
    vagrant ssh -c 'cd share && for c in registry logger database \
      cache controller builder router; \
      do cd $c && sudo systemctl enable $(pwd)/systemd/* && cd ..; done'
    ln -s '/home/core/share/registry/systemd/deis-registry.service' \
      '/etc/systemd/system/multi-user.target.wants/deis-registry.service'
    ...
    Apr 15 18:53:23 deis sh[9101]: 2014-04-15 12:53:23 [149] [INFO] Booting worker with pid: 149
    Apr 15 18:53:24 deis sh[9101]: [2014-04-15 12:53:24,842: INFO/MainProcess] mingle: all alone
    Apr 15 18:53:24 deis sh[9101]: [2014-04-15 12:53:24,852: WARNING/MainProcess] celery@121f56ff9ae5 ready.

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

Once the user is registered, add your SSH key for ``git push``
access using:

.. code-block:: console

    $ deis keys:add
    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done


Your local development environment is running! Follow the
rest of the :ref:`Developer Guide <developer>` to deploy your first application.

Test Your Changes
-----------------
In the single-node Vagrant environment, testing your changes to Deis itself
is easy:

    - Make changes to the code in one of the component subdirectories, such
      as ``controller/``
    - run ``make -C controller/ build run``
    - Test your changes with ``make -C controller/ test`` and interactively
      with the Deis client

Useful Commands
---------------

Once your controller is running, here are some helpful commands.

Tail Logs
`````````

.. code-block:: console

    $ vagrant ssh -c 'sudo docker logs --follow=true deis-controller'

Restart Services
````````````````

.. code-block:: console

    $ vagrant ssh -c 'sudo restart deis-controller'

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

Creating a pull request on the Deis repository also runs a Travis CI build to
ensure the pull request doesn't break any tests or reduce code coverage.


.. _`easy-fix`: https://github.com/deis/deis/issues?labels=easy-fix&state=open
.. _`Vagrant`: http://www.vagrantup.com/
.. _`VirtualBox`: https://www.virtualbox.org/
.. _`fork the Deis repository`: https://github.com/deis/deis/fork
.. _`pull request`: https://github.com/deis/deis/pulls
