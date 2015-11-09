:description: How to hack on Deis including setup instructions

.. _hacking:

Development Environment
=======================

DigitalOcean_ is the recommended development environment for Deis project
contributors. :ref:`Provision a new DO cluster <deis_on_digitalocean>` and then
continue to follow the instructions below to get started hacking.

.. _digitalocean_credit:

.. important::

    Are you a new contributor to Deis? Your first `Pull Request`_ could earn you
    credit at DigitalOcean_! Submit your changes and then email
    deis@engineyard.com. When your PR is merged, the maintainer team will
    send you a DigitalOcean credit based on the value of your contribution.

This document is for developers who are interested in working directly on the
Deis codebase. In this guide, we walk you through the process of setting up
a local development environment. While there are many ways to set up your
Deis environment, this document covers a specific setup:

- Developing on **Mac OSX** or **Linux**
- Managing virtualization with **Vagrant/Virtualbox**
- Hosting a docker registry with **docker-machine** (Mac)

We try to make it simple to hack on Deis. However, there are necessarily several moving
pieces and some setup required. We welcome any suggestions for automating or simplifying
this process.

If you're just getting into the Deis codebase, look for GitHub issues with the label
`easy-fix`_. These are more straightforward or low-risk issues and are a great way to
become more familiar with Deis.

Prerequisites
-------------

You can develop on any supported platform including your laptop, cloud providers or
on bare metal.  We strongly recommend a minimum 3-node cluster. We strongly
suggest using Vagrant and VirtualBox for your virtualization layer during
development.

At a glance, you will need:

- Python 2.7 or later (with ``pip``)
- virtualenv (``sudo pip install virtualenv``)
- Go 1.5 or later, with support for compiling to ``linux/amd64``
- Godep (https://github.com/tools/godep)
- VirtualBox latest
- Vagrant 1.5 or later
- On Mac, you will also want
  - Docker Machine (http://docs.docker.com/machine/install-machine/)

Additionally, you may need:
- shellcheck (https://github.com/koalaman/shellcheck)
- golint (https://github.com/golang/lint)

In most cases, you should simply install according to the instructions. There
are a few special cases, though. We cover these below.

Configuring Go
``````````````

If your local workstation does not support the linux/amd64 target environment, you will
have to install Go from source with cross-compile support for that environment. This is
because some of the components are built on your local machine and then injected into a
docker container.

Homebrew users can just install with cross compiling support:

.. code-block:: console

    $ brew install go --with-cc-common

It is also straightforward to build Go from source:

.. code-block:: console

    $ sudo su
    $ curl -sSL https://golang.org/dl/go1.5.src.tar.gz | tar -v -C /usr/local -xz
    $ cd /usr/local/go/src
    $ # compile Go for our default platform first, then add cross-compile support
    $ ./make.bash --no-clean
    $ GOOS=linux GOARCH=amd64 ./make.bash --no-clean

Once you can compile to ``linux/amd64``, you should be able to compile Deis'
components as normal.

Configuring Docker Machine (Mac)
````````````````````````````````

Deis needs a Docker registry running independently of the Deis cluster. On
OS X, you will need Docker Machine (http://docs.docker.com/machine/install-machine/)
to run the registry inside of a VirtualBox image.

.. note::

    Previously, Deis used boot2docker to run the registry. However, Docker has
    deprecated boot2docker in favor of Docker Machine.

Install Docker Machine according to the normal installation instructions. Then
create a new image for hosting your Deis Docker registry:

.. code-block:: console

    $ docker-machine create --driver virtualbox --virtualbox-disk-size=100000 \
    --engine-insecure-registry=192.168.0.0/16 deis-registry

This will create a new virtual machine named `deis-registry` that will take
up as much as 100,000 MB of disk space. Registries tend to be large, so
allocating a big disk is a good idea.

Once the deis-registry machine exists, source its values into your environment
so your docker client knows how to use the new machine.

.. code-block:: console

    $ eval "$(docker-machine env deis-registry)"

.. note::

    Because the registry that we create will not have a valid SSL certificate,
    we run the local registry as an insecure (HTTP, not HTTPS) registry. Each
    time Docker Machine reboots, the registry will get a new IP address
    somewhere in the 192.168.0.0/16 range. We must declare that explicitly when
    configuring Docker Machine.

At this point, our `deis-registry` VM can now serve as a registry for Deis'
Docker images. Later we will return to this.

Fork the Deis Repository
------------------------
Once the prerequisites have been met, we can begin to work with Deis.

To get Deis running for development, first `fork the Deis repository`_,
then clone your fork of the repository. Since Deis is predominantly written
in Go, the best place to put it is in ``$GOPATH/src/github.com/deis/``

.. code-block:: console

    $ mkdir -p  $GOPATH/src/github.com/deis
    $ cd $GOPATH/src/github.com/deis
    $ git clone git@github.com:<username>/deis.git
    $ cd deis

.. note::

    By checking out the forked copy into the namespace ``github.com/deis/deis``,
    we are tricking the Go toolchain into seeing our fork as the "official"
    Deis tree.

If you are going to be issuing pull requests and working with official Deis
repository, we suggest configuring Git accordingly. There are various strategies
for doing this, but the `most common`_ is to add an ``upstream`` remote:

.. code-block:: console

    $ git remote add upstream https://github.com/deis/deis.git

For the sake of simplicity, you may want to point an environment variable to
your Deis code:

.. code-block:: console

    export DEIS=$GOPATH/src/github.com/deis/deis

Throughout the rest of this document, ``$DEIS`` refers to that location.

Alternative: Forking with a Pushurl
```````````````````````````````````
A number of Deis developers prefer to pull directly from ``deis/deis``, but
push to ``<username>/deis``. If that workflow suits you better, you can set it
up this way:

.. code-block:: console

    $ git clone git@github.com:deis/deis.git
    $ cd deis
    $ git config remote.origin.pushurl git@github.com:<username>/deis.git

In this setup, fetching and pulling code will work directly with the upstream
repository, while pushing code will send changes to your fork. This makes it
easy to stay up to date, but also make changes and then issue pull requests.

Build deisctl
-------------

``deisctl`` is used for interacting with the Deis cluster. While you can use an
existing ``deisctl`` binary, we recommend that developers build it from source.

.. code-block:: console

  $ cd $DEIS/deisctl
  $ make build
  $ make install  # optionally

This will build just the ``deisctl`` portion of Deis. Running ``make install`` will
install the ``deisctl`` command in ``$GOPATH/bin/deisctl``.

You can verify that ``deisctl`` is correctly built and installed by running
``deisctl -h``. That should print the help text and exit.

Configure SSH Tunneling for Deisctl
-----------------------------------

To connect to the cluster using ``deisctl``, you must add the private key to ``ssh-agent``.
For example, when using Vagrant:

.. code-block:: console

    $ ssh-add ~/.vagrant.d/insecure_private_key

Set ``DEISCTL_TUNNEL`` so the ``deisctl`` client on your workstation can connect to
one of the hosts in your cluster:

.. code-block:: console

    $ export DEISCTL_TUNNEL=172.17.8.100

.. note::

  A number of times during this setup, tools will suggest that you export various
  environment variables. You may find it convenient to store these in your shell's
  RC file (`~/.bashrc` or `~/.zshrc`).

Install the Deis Client
-----------------------

The ``deis`` client is also written in Go. Your Deis client should match your server's
version. Like ``deisctl``, we recommend that developers build ``deis`` from source:

.. code-block:: console

    $ cd $DEIS/client
    $ make build
    $ make install  # optionally
    $ ./deis
    Usage: deis <command> [<args>...]


Start Up a Development Cluster
------------------------------

Our host system is now configured for controlling a Deis cluster. The next
thing to do is begin standing up a development cluster.

When developing locally, we want deisctl to check our local unit files so that
any changes are reflected in our Deis cluster. The easiest way to do this is
to set an environment variable telling deisctl where to look. Assuming
the variable ``$DEIS`` points to the location if the deis source code, we want
something like this:

.. code-block:: console

    export DEISCTL_UNITS=$DEIS/deisctl/units

To start up and configure a local vagrant cluster for development, you can use
the ``dev-cluster`` target.

.. code-block:: console

    $ make dev-cluster

This may take a while to run the first time. At the end of the process, you
will be prompted to run ``deis start platform``. Hold off on that task for now.
We will come back to it later.

To verify that the cluster is running, you should be able to connect
to the nodes on your Deis cluster:

.. code-block:: console

    $ vagrant status
    Current machine states:

    deis-01               running (virtualbox)
    deis-02               running (virtualbox)
    deis-03               running (virtualbox)

    $ vagrant ssh deis-01
    Last login: Tue Jun  2 18:26:30 2015 from 10.0.2.2
     * *    *   *****    ddddd   eeeeeee iiiiiii   ssss
    *   *  * *  *   *     d   d   e    e    i     s    s
     * *  ***** *****     d    d  e         i    s
    *****  * *    *       d     d e         i     s
    *   * *   *  * *      d     d eee       i      sss
    *****  * *  *****     d     d e         i         s
      *   *****  * *      d    d  e         i          s
     * *  *   * *   *     d   d   e    e    i    s    s
    ***** *****  * *     ddddd   eeeeeee iiiiiii  ssss

    Welcome to Deis			Powered by CoreOS

With a dev cluster now running, we are ready to set up a local Docker registry.

Configure a Docker Registry
---------------------------

The development workflow requires Docker Registry set at the ``DEV_REGISTRY``
environment variable.  If you're developing locally you can use the ``dev-registry``
target to spin up a quick, disposable registry inside a Docker container.
The target ``dev-registry`` prints the registry's address and port when using ``docker-machine``;
otherwise, use your host's IP address as returned by ``ifconfig`` with port 5000 for ``DEV_REGISTRY``.

.. code-block:: console

    $ make dev-registry

    To configure the registry for local Deis development:
        export DEV_REGISTRY=192.168.59.103:5000

It is important that you export the ``DEV_REGISTRY`` variable as instructed.

If you are developing elsewhere, you must set up a registry yourself.
Make sure it meets the following requirements:

 #. You can push Docker images from your workstation
 #. Hosts in the cluster can pull images with the same URL

.. note::

    If the development registry is insecure and has an IP address in a range other than ``10.0.0.0/8``,
    ``172.16.0.0/12``, or ``192.168.0.0/16``, you'll have to modify ``contrib/coreos/user-data.example``
    and whitelist your development registry so the daemons can pull your custom components.

Initial Platform Build
----------------------

The full environment is prepared. You can now build Deis from source code and
then run the platform.

We'll do three steps together:

- Build the source (``make build``)
- Update our local cluster with a dev release (``make dev-release``)
- Start the platform (``deisctl start platform``)

Conveniently, we can accomplish all three in one step:

.. code-block:: console

    $ make deploy


Running ``deisctl list`` should display all of the services that your Deis
cluster is currently running.

You can now use your Deis cluster in all of the usual ways.

At this point, you are running Deis from the code in your Git clone. But since
rebuilding like this is time consuming, Deis has a simplified developer
workflow more suited to daily development.

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

.. note::

   You can export the ``DEIS_STATELESS=True`` environment variable to skip all
   store components when using the root ``Makefile``. Useful when working
   on a stateless platform (:ref:`running-deis-without-ceph`).

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

    $ deisctl list             # determine which host runs the controller
    $ ssh core@<host>          # SSH into the controller host
    $ nse deis-controller      # inject yourself into the container
    $ cd /app                  # change into the django project root
    $ ./manage.py shell        # get a django shell

Have commands other Deis developers might find useful? Send us a PR!

Pull Requests
-------------

Please read :ref:`standards`. It contains a checklist of things you should do
when proposing a change to Deis.

.. _DigitalOcean: https://www.digitalocean.com/
.. _`Pull Request`: https://github.com/deis/deis/pulls
.. _`easy-fix`: https://github.com/deis/deis/issues?labels=easy-fix&state=open
.. _`deisctl`: https://github.com/deis/deis/tree/master/deisctl
.. _`fork the Deis repository`: https://github.com/deis/deis/fork
.. _`Python 2.7`: https://www.python.org/downloads/release/python-279/
.. _`running the tests`: https://github.com/deis/deis/tree/master/tests#readme
.. _`pull request`: https://github.com/deis/deis/pulls
.. _`most common`: https://help.github.com/articles/fork-a-repo/
