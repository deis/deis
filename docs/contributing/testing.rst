:title: Testing Deis
:description: How to test Deis

.. _testing:

Testing Deis
============

Deis is a distributed system with many moving parts, which makes it of paramount
importance to test every change thoroughly.

Deis is also a set of components that correspond to directories in the source
code repository. Most components are Docker containers, two are command-line
clients, and one contains the documentation. Components have source-code level
`unit tests`_ and black-box type `functional tests`_. `integration tests`_
verify the behavior of the components together as a system.

GitHub pull requests for Deis are tested automatically by a Jenkins
`continuous integration`_ (CI) system at https://ci.deis.io. Contributors should
run the same tests locally before proposing any changes to the Deis codebase.


Set Up the Environment
----------------------

To run all tests, you will need:

- Vagrant 1.6.5 or later
- VirtualBox 4.3 or later
- Docker 1.3.0
- PostgreSQL server

The tests assume that you have Deis' `source code`_ in your ``$GOPATH``:

.. code-block:: console

    $ go get -u -v github.com/deis/deis
    $ cd $GOPATH/src/github.com/deis/deis

Start a Docker Registry
^^^^^^^^^^^^^^^^^^^^^^^

Deis' functional tests build Docker images and test them locally. The images are
then pushed to a `Docker registry`_ so that integration tests can test them as
binary artifacts--just as a real-world provisioning of Deis pulls images from
the Docker Hub.

If you don't have a Docker registry already accessible for your testing or for
continuous deployment, start one locally.

.. code-block:: console

    $ make dev-registry
    registry

    To use local registry for Deis development:
        export DEV_REGISTRY=192.168.59.103:5000

.. important::

    The functional tests also use several mock or example containers:
    **deis/test-etcd**, **deis/test-postgresql**, and **deis/mock-store**.
    These are built locally during a test run.

Run the Tests
-------------

The unit and functional tests for each component are in their respective
directories. The integration tests, scripts, and supporting go packages are in
the ``tests/`` directory in the project root.

Scripts in the ``tests/bin`` directory are the best place to start. These test
individual pieces of Deis, then bring up a Vagrant cluster and test all of them
as a system. They call ``tests/bin/test-setup.sh`` to test for important
environment variables and will exit with a helpful message if any are missing.

The ``test-setup.sh`` script also prepares the testing environment, as well as
tears it down after testing is complete. If there is a test failure, the script
collects verbose component logs, compresses them, and places them in ``$HOME``.
If `s3cmd`_ is installed and configured on the test machine, the script will
instead upload the logs to Amazon S3. This is how the Jenkins CI infrastructure
is configured, so that contributors have access to the logs to see how their
PR failed.

test-integration.sh
^^^^^^^^^^^^^^^^^^^

- runs documentation tests
- builds Docker images tagged with ``$BUILD_TAG``
- runs unit and functional tests
- creates a 3-node Vagrant CoreOS cluster
- pushes the Docker images to a registry
- provisions the cluster for Deis with the registry images
- runs all integration tests
- takes roughly an hour

.. code-block:: console

    $ ./tests/bin/test-integration.sh

    >>> Preparing test environment <<<

    DEIS_ROOT=/Users/matt/Projects/src/github.com/deis/deis
    DEIS_TEST_APP=example-dockerfile-http
    ...
    >>> Running integration suite <<<

    make -C tests/ test-full
    ...
    >>> Test run complete <<<

test-smoke.sh
^^^^^^^^^^^^^

- runs documentation tests
- builds Docker images tagged with ``$BUILD_TAG``
- runs unit and functional tests
- creates a 3-node Vagrant CoreOS cluster
- pushes the Docker images to a registry
- provisions the cluster for Deis with the registry images
- runs a "smoke test" that pushes and scales an app
- takes roughly 45 minutes

test-latest.sh
^^^^^^^^^^^^^^

- installs the latest ``deis`` and ``deisctl`` client releases
- creates a 3-node Vagrant CoreOS cluster
- provisions the cluster for Deis with latest release images
- runs a "smoke test" that pushes and scales an app
- takes roughly 30 minutes

Run Specific Tests
^^^^^^^^^^^^^^^^^^

Run the tests for a single component this way:

.. code-block:: console

    $ make -C logger test             # unit + functional
    $ make -C controller test-unit
    $ make -C router test-functional


Customize Test Runs
-------------------

The file ``tests/bin/test-setup.sh`` is the best reference to environment
variables that can affect the tests' behavior. Here are some important ones:

- ``HOST_IPADDR`` - address on which Docker containers can communicate for the
  functional tests, probably the host's IP or the one assigned to `Docker Machine`_.
- ``DEIS_TEST_APP`` - name of the `Deis example app`_ to use, which is cloned
  from GitHub (default: ``example-go``)
- ``DEIS_TEST_AUTH_KEY`` - SSH key used to register with the Deis controller
  (default: ``~/.ssh/deis``)
- ``DEIS_TEST_SSH_KEY`` - SSH key used to login to the controller machine
  (default: ``~/.vagrant.d/insecure_private_key``)
- ``DEIS_TEST_DOMAIN`` - the domain to use for testing
  (default: ``local3.deisapp.com``)


.. _`unit tests`: http://en.wikipedia.org/wiki/Unit_testing
.. _`functional tests`: http://en.wikipedia.org/wiki/Functional_testing
.. _`integration tests`: http://en.wikipedia.org/wiki/Integration_testing
.. _`continuous integration`: http://en.wikipedia.org/wiki/Continuous_integration
.. _`Docker Machine`: http://docs.docker.com/machine/install-machine/
.. _`source code`: https://github.com/deis/deis
.. _`Docker registry`: https://github.com/docker/docker-registry
.. _`Deis example app`: https://github.com/deis?query=example-
.. _`s3cmd`: http://s3tools.org/s3cmd
