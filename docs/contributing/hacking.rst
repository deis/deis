:description: How to hack on Deis including setup instructions

.. _hacking:

Hacking on Deis
===============

We try to make it simple to hack on Deis. However, there are necessarily several moving
pieces and some setup required. We welcome any suggestions for automating or simplifying
this process.

If you're just getting into the Deis codebase, look for GitHub issues with the label
`easy-fix`_. These are more straightforward or low-risk issues and are a great way to
become more familiar with Deis.

Prerequisites
-------------

You can develop on any supported platform including your laptop, cloud providers or
on bare metal.  We strongly recommend a minimum 3-node cluster.

The development workflow requires a Docker Registry that is accessible to you
(the developer) and to all of the hosts in your cluster.

You will also need a `deisctl`_ client to update images and restart components.

Fork the Repository
-------------------

To get Deis running for development, first `fork the Deis repository`_,
then clone your fork of the repository:

.. code-block:: console

	$ git clone git@github.com:<username>/deis.git
	$ cd deis
	$ export DEIS_DIR=`pwd`  # to use in future commands

Install the Client
------------------

In a development environment you'll want to use the latest version of the client. Install
its dependencies by using the Makefile and symlinking ``client/deis.py`` to ``deis`` on
your local workstation.

.. code-block:: console

    $ cd $DEIS_DIR/client
    $ make install
    $ sudo ln -fs $DEIS_DIR/client/deis.py /usr/local/bin/deis
    $ deis
    Usage: deis <command> [<args>...]

Configure SSH Tunneling
-----------------------

To connect to the cluster using ``deisctl``, you must add the private key to ``ssh-agent``.
For example, when using Vagrant:

.. code-block:: console

    $ ssh-add ~/.vagrant.d/insecure_private_key

Set ``DEISCTL_TUNNEL`` so the ``deisctl`` client on your workstation can connect to
one of the hosts in your cluster:

.. code-block:: console

    $ export DEISCTL_TUNNEL=172.17.8.100

Test connectivity using ``deisctl list``:

.. code-block:: console

    $ deisctl list

Configure a Docker Registry
---------------------------

The development workflow requires Docker Registry set at the ``DEV_REGISTRY``
environment variable.  If you're developing locally you can use the ``dev-registry``
target to spin up a quick, disposable registry inside a Docker container.

.. code-block:: console

    $ make dev-registry

    To configure the registry for local Deis development:
        export DEV_REGISTRY=192.168.59.103:5000

.. note::

	For Docker 1.3.1 and later, ``docker push`` to this development registry may fail
	without SSL certificate support. Restart docker with an ``--insecure-registry`` flag.

	For ``boot2docker`` 1.3.1 for example, add
	``EXTRA_ARGS="--insecure-registry 192.168.59.103:5000"`` to
	/var/lib/boot2docker/profile and restart docker with ``sudo /etc/init.d/docker restart``.

If you are developing elsewhere, you must set up a registry yourself.
Make sure it meets the following requirements:

 #. You can push Docker images from your workstation
 #. Hosts in the cluster can pull images with the same URL

.. note::

    If the development registry is insecure and has an IP address in a range other than ``10.0.0.0/8``,
    ``172.16.0.0/12``, or ``192.168.0.0/16``, you'll have to modify ``contrib/coreos/user-data.example``
    and whitelist your development registry so the daemons can pull your custom components.

Development Workflow
--------------------

Deis includes ``Makefile`` targets designed to simplify the development workflow.
This workflow is typically:

  #. Update source code and commit your changes using ``git``
  #. Use ``make -C <component> build`` to build a new Docker image
  #. Use ``make -C <component> dev-release`` to push a snapshot release
  #. Use ``make -C <component> restart`` to restart the component

This can be shortened to a one-liner using the ``deploy`` target:

.. code-block:: console

    $ make -C controller deploy

You can also use the same tasks on the root ``Makefile`` to operate on all
components at once.  For example, ``make deploy`` will build, dev-release,
and restart all components on the cluster.

.. important::

   In order to cut a dev-release, you must commit changes using ``git`` to increment
   the SHA used when tagging Docker images

Test Your Changes
-----------------

Deis ships with a comprehensive suite of automated tests, most written in Go.
See :ref:`testing` for instructions on running the tests.

Useful Commands
---------------

Once your controller is running, here are some helpful commands.

Tail Logs
`````````

.. code-block:: console

    $ deisctl journal controller

Rebuild Services from Source
````````````````````````````

.. code-block:: console

    $ make -C controller build push restart

Restart Services
````````````````

.. code-block:: console

    $ make -C controller restart

Django Shell
````````````

.. code-block:: console

    $ deisctl ssh controller   # SSH into the controller
    $ nse deis-controller      # inject yourself into the container
    $ cd /app                  # change into the django project root
    $ ./manage.py shell        # get a django shell

Have commands other Deis developers might find useful? Send us a PR!

Pull Requests
-------------

Please read :ref:`standards`. It contains a checklist of things you should do
when proposing a change to Deis.

.. _`easy-fix`: https://github.com/deis/deis/issues?labels=easy-fix&state=open
.. _`deisctl`: https://github.com/deis/deis/tree/master/deisctl
.. _`fork the Deis repository`: https://github.com/deis/deis/fork
.. _`running the tests`: https://github.com/deis/deis/tree/master/tests#readme
.. _`pull request`: https://github.com/deis/deis/pulls
