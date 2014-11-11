:title: System Requirements
:description: System requirements for provisioning Deis.

.. _system-requirements:

System Requirements
===================

When deploying Deis, it's important to provision machines with adequate resources. Deis is a
highly-available distributed system, which means that Deis components and your deployed applications
will move around the cluster onto healthy hosts as hosts leave the cluster for various reasons
(failures, reboots, autoscalers, etc.). Because of this, you should have ample spare resources on
any machine in your cluster to withstand the additional load of running services for failed machines.

Machines must have:

* At least 4GB of RAM (Deis uses 2 - 2.5GB, plus room for failover and deployed applications)
* At least 40GB of hard disk space

Running smaller machines will likely result in increased system load and has been known to result
in component failures, issues with etcd/fleet, and other problems.

If running multiple (at least 3) machines of an adequate size is unreasonable, it is recommended to
investigate the `Dokku`_ project instead. Dokku is `sponsored`_ by Deis and is ideal for environments
where a highly-available distributed system is not necessary (i.e. local development, testing, etc.).

.. _`dokku`: https://github.com/progrium/dokku
.. _`sponsored`: http://deis.io/deis-sponsors-dokku/
