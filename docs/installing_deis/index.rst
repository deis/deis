:title: Installing Deis
:description: Step-by-step guide for operations engineers setting up a private PaaS using Deis.

.. _installing_deis:
.. _provision-controller:

Installing Deis
================

The `controller` is the brains of a Deis platform. Provisioning a Deis
controller is a matter of creating one or more :ref:`concepts_coreos`
machines and installing a few necessary *systemd* units to manage
Docker containers.

Anywhere you can run CoreOS, you can run Deis, including most cloud
providers, virtual machines, and bare metal. See the
`CoreOS documentation`_ for more information on how to get set up
with CoreOS.

:Release: |version|
:Date: |today|

.. toctree::

    digitalocean
    aws
    vagrant
    gce
    rackspace
    baremetal


.. _`CoreOS Documentation`: https://coreos.com/docs/
