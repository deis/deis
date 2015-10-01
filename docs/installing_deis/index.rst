:title: Installing Deis
:description: Step-by-step guide for operations engineers setting up Deis.

.. _installing_deis:

Installing Deis
===============

Installing Deis consists of provisioning three or more :ref:`concepts_coreos`
machines and using :ref:`deisctl <install_deisctl>` to set up and start
the core components.

Anywhere you can run CoreOS, you can run Deis, including most cloud providers,
virtual machines, and bare metal. See the `CoreOS documentation`_ for more
information on how to get set up with CoreOS.

:Release: |version|
:Date: |today|

.. toctree::

    quick-start
    system-requirements
    aws
    baremetal
    digitalocean
    gce
    linode
    azure
    openstack
    vagrant
    install-deisctl
    install-platform

.. _`CoreOS Documentation`: https://coreos.com/docs/
