:description: Developer setup instructions for contributing to the Deis project.
:keywords: deis, documentation, contributing, developer, setup, chef, knife

.. _devsetup:

Developer Setup
===============

Thank you for contributing to Deis! We have tried to make it simple
to work on Deis, but as an open PaaS there are necessarily several
moving pieces and some setup required. We welcome any suggestions for
automating or simplifying this process.


Prerequisites
-------------

We assume you have a modern UNIX-like environment, such as Linux or
Mac OS X. Many of the tools listed below will probably be present
already, but install those that are not:

- git
- GNU Make
- `Python`_ 2.7.x
- `Ruby`_ 1.9 (for knife.rb)
- `Amazon EC2 API Tools`_

To contribute code back to Deis, you must also have a GitHub.com account
in order to create a pull request.


.. _first_things_first:

First Things First
------------------

To work on Deis itself, first `fork the Deis repository`_ at GitHub.com.
Then clone *your* repository for local development:

.. code-block:: console

	$ git clone https://github.com/<username>/deis.git
	$ cd deis


Don't clone the official repository, but **do** complete all the other steps
outlined in :ref:`installation`. Ensure your CLI and controller are functional
before moving on.

Near the end of running the ``contrib/aws/provision-ec2-controller.sh`` script,
you will see output similar to this::

	Instance ID: i-38ad000c
	Flavor: m1.large
	Image: ami-b55ac885
	Region: us-west-2
	SSH Key: deis-controller
	Public DNS Name: ec2-198-51-100-36.us-west-2.compute.amazonaws.com
	Public IP Address: 198.51.100.36
	Run List: recipe[deis], recipe[deis::server], recipe[deis::gitosis],
	...

Note the **Public DNS Name** value (**ec2-198-51-100-36.us-west-2.compute.amazonaws.com**
in this example). This is the Amazon EC2 instance that runs
your Deis controller software, and ultimately this is where you will test
any changes you make to the Deis codebase.


Make a Virtualenv
-----------------

To keep Deis` requirements separate from other development you may do,
it's preferable to create a **virtual environment** for python.

.. code-block:: console

	$ virtualenv venv --prompt='(deis)'
	New python executable in venv/bin/python
	Installing Setuptools.................................done.
	Installing Pip........................................done.
	$ source venv/bin/activate
	(deis)$ pip install -r dev_requirements.txt --use-mirrors
	Downloading/unpacking boto==2.13.3 (from -r requirements.txt (line 2))
  	Downloading boto-2.13.3.tar.gz (1.1MB): 1.1MB downloaded
  	Running setup.py egg_info for package boto
    Downloading/unpacking celery==3.0.22 (from -r requirements.txt (line 3))
	...
	Successfully installed boto celery ...
	Cleaning up...
	(deis)$

Make sure you install the requirements in the dev_requirements.txt file,
which contains several additions over the runtime requirements.txt file.
Please see the `virtualenv documentation`_ for more details on python virtual
environments.


Modify Code and Test
--------------------

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


Test on Your Controller
-----------------------

Since you completed the steps outlined in :ref:`first_things_first`, you have
a working Deis controller. Start a remote shell on the controller as the
"ubuntu" user:

.. code-block:: console

	$ ssh -i $HOME/.ssh/deis-controller ubuntu@ec2-198-51-100-36.us-west-2.compute.amazonaws.com
	Welcome to Ubuntu 12.04.2 LTS (GNU/Linux 3.8.0-26-generic x86_64)

	ubuntu@ip-198-51-100-36:~$
	$ status deis-server
	deis-server start/running, process 12040
	$ cd /opt/deis
	$ sudo -u deis mv controller controller.opdemand
	$ # clone my fork of the official Deis repo
	$ sudo -u deis git clone https://github.com/<username>/deis.git controller
	Cloning into 'controller'...
	remote: Counting objects: 2067, done.
	remote: Compressing objects: 100% (1007/1007), done.
	remote: Total 2067 (delta 951), reused 2064 (delta 948)
	Receiving objects: 100% (2067/2067), 1.01 MiB | 924 KiB/s, done.
	Resolving deltas: 100% (951/951), done.
	$ ls
	build  controller  controller.opdemand  gitosis  prevent-apt-update
	$ cd controller
	$ sudo -u deis cp controller.opdemand/deis/local_settings.py controller/deis/
	$ sudo -u deis cp controller.opdemand/.secret_key controller/
	$ sudo -u deis cp -r controller.opdemand/venv controller/
	$ sudo restart deis-server
	deis-server start/running, process 12901

You have now restarted the Deis controller from your fork of the codebase.

Testing now involves exercising the relevant code paths in full round-trip
mode by using the ``deis`` client on your workstation. You can get detailed
output by editing the deis/local_settings.py file and set DEBUG=True,
restarting the server, and watching logs:

.. code-block:: console

	$ sudo restart deis-server
	deis-server start/running, process 14074
	$ tail -f /var/log/deis/*.log
	==> /var/log/deis/access.log <==

	==> /var/log/deis/celeryd.log <==
	[2013-08-13 16:59:33,426: WARNING/MainProcess] celery@ip-198-51-100-36 ready.
	[2013-08-13 16:59:33,451: INFO/MainProcess] consumer: Connected to amqp://guest@127.0.0.1:5672//.

	==> /var/log/deis/server.log <==
	2013-08-13 23:29:09 [14019] [INFO] Handling signal: term
	2013-08-13 23:29:09 [14019] [INFO] Shutting down: Master
	2013-08-13 23:29:12 [14074] [INFO] Starting gunicorn 17.5
	2013-08-13 23:29:12 [14074] [DEBUG] Arbiter booted
	2013-08-13 23:29:12 [14074] [INFO] Listening at: http://0.0.0.0:8000 (14074)
	2013-08-13 23:29:12 [14074] [INFO] Using worker: gevent
	2013-08-13 23:29:12 [14079] [INFO] Booting worker with pid: 14079


Make a Pull Request
-------------------

Please create a GitHub `pull request`_ for any code changes that will benefit Deis users
in general. This workflow helps changesets map well to discrete features.

Creating a pull request on the Deis repository also runs a `Travis CI`_ build to
ensure the pull request doesn't break any tests or reduce code coverage.


Clone the Deis Cookbook
-----------------------

If you want to modify Deis' Chef recipes, you should also clone the `deis-cookbook`_
repository:

.. code-block:: console

	$ git clone -q https://github.com/opdemand/deis-cookbook.git

Please see `deis-cookbook`_ for information about contributing Chef code to Deis.


.. _`virtualenv documentation`: http://www.virtualenv.org/en/latest/
.. _`Python`: http://python.org/
.. _`Ruby`: http://ruby-lang.org/
.. _`Amazon EC2 API Tools`: http://aws.amazon.com/developertools/Amazon-EC2/351
.. _`Knife EC2 plugin`: https://github.com/opscode/knife-ec2
.. _`fork the Deis repository`: https://github.com/opdemand/deis/fork
.. _`deis-cookbook`: https://github.com/opdemand/deis-cookbook
.. _`pull request`: https://github.com/opdemand/deis/pulls
.. _`Travis CI`: https://travis-ci.org/
