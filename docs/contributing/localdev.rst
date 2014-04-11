:description: Local development setup instructions for contributing to the Deis project.

.. _localdev:

Local Development
=================
We have tried to make it simple to hack on Deis, but as an open PaaS, there are
necessarily several moving pieces and some setup required. We welcome any suggestions for
automating or simplifying this process.

Prerequisites
-------------
We strongly recommend using `Vagrant`_ with `VirtualBox`_ so you can develop inside a set
of isolated virtual machines. You will need:

 * Vagrant 1.3.5+
 * VirtualBox 4.2+

Fork the Repository
-------------------
To get Deis running locally, first `fork the Deis repository`_, then clone your fork of
the repository for local development:

.. code-block:: console

	$ git clone git@github.com:<username>/deis.git
	$ cd deis
	$ export DEIS_DIR=`pwd`  # to use in future commands

Provision the Controller
------------------------
.. code-block:: console

    $ vagrant up

Yes, really. That's it.

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

    $ deis register http://local.deisapp.com:8000
    username: myuser
    password:
    password (confirm):
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

.. note::

    As of v0.5.1, the proxy was removed for Deis platform services. It has yet to be added
    back in. See `issue 535`_ for more details.

    As a workaround, use the following:

    :code:`deis register http://local.deisapp.com:8000`

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


.. _`Vagrant`: http://www.vagrantup.com/
.. _`VirtualBox`: https://www.virtualbox.org/
.. _`fork the Deis repository`: https://github.com/opdemand/deis/fork
.. _`pull request`: https://github.com/opdemand/deis/pulls
.. _`issue 535`: https://github.com/opdemand/deis/issues/535
