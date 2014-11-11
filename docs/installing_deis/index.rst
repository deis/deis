:title: Installing Deis
:description: Step-by-step guide for operations engineers setting up a private PaaS using Deis.

.. _installing_deis:
.. _provision-controller:

Installing Deis
===============

Provisioning Deis is a matter of creating one or more :ref:`concepts_coreos`
machines and using :ref:`install_deisctl` to install and start Deis.

Anywhere you can run CoreOS, you can run Deis, including most cloud
providers, virtual machines, and bare metal. See the
`CoreOS documentation`_ for more information on how to get set up
with CoreOS.

:Release: |version|
:Date: |today|

.. toctree::

    quick-start
    system-requirements
    aws
    digitalocean
    gce
    rackspace
    vagrant
    baremetal
    install-deisctl
    install-platform

.. _`CoreOS Documentation`: https://coreos.com/docs/
