:title: Security considerations
:description: Security considerations for Deis.

.. _security_considerations:

Security considerations
========================

A major goal of Deis is to be operationally secure and trusted by operations engineers in every deployed
environment. To reach that goal, several security concerns need to be addressed.

Access to etcd
--------------
Since all Deis configuration settings are stored in etcd (including passwords, keys, etc.), any access
to the etcd cluster compromises the security of the entire Deis installation. The various provision
scripts configure the etcd daemon to only listen on the private network interface, but any host or
container with access to the private network has full access to etcd. This also includes deployed
application containers, which cannot be trusted.

The planned approach is to configure iptables on the machines to prevent unauthorized access from
containers. Some requirements include:

* Containers must be able to access the outside world
* Containers must be able to access other containers
* Containers cannot access the CoreOS host (SSH, etcd, etc)

Further discussion about this approach is appreciated in GitHub issue `#986`_.

Separate runtime environments
-----------------------------
A design goal of Deis is to become scheduler agnostic. While the core Deis components (builder,
cache, controller, database) will remain on a CoreOS cluster running etcd, application containers
will be scheduled to an external cluster. The registry, logger, and router components will live on
each scheduling cluster. This enables Deis to use other scheduling algorithms, and it also introduces
a clean network segregation. This will alleviate the concern that application containers have access
to the core Deis etcd installation, but access policies between clusters will need to be introduced.

.. _`#986`: https://github.com/deis/deis/issues/986
