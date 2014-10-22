:title: Build Deis Components
:description: Building from source

.. _build-components:

Build Deis Components
=====================

Pre-requisites
--------------

Install golang in your development environment. To build ``deisctl``, you will need to install these go libraries:

.. code-block:: console

    $ go get github.com/tools/godep
    $ go get github.com/golang/lint
    $ github.com/deis/deis/deisctl/backend/fleet
    $ cd deisctl && make build

Build
-----

.. code-block:: console

    $ make build dev-release

Starting Components
-------------------

.. code-block:: console

    $ export DEISCTL_UNITS="$GOPATH/src/github.com/deis/deis/deisctl/units"
    $ deisctl config platform set domain="example.com"
    $ deisctl install platform
    $ deisctl start platform
