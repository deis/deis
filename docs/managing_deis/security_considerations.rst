:title: Security considerations
:description: Security considerations for Deis.

.. _security_considerations:

Security considerations
========================

.. important::

    Deis is not suitable for multi-tenant environments
    or hosting untrusted code.

A major goal of Deis is to be operationally secure and trusted by operations engineers in every deployed
environment. There are, however, two notable security-related considerations to be aware of
when deploying Deis.


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

Application runtime segregation
-------------------------------
Users of Deis often want to deploy their applications to separate environments
(commonly: development, staging, and production). Typically, physical network isolation isn't
the goal, but rather segregation of application environments - if a development app goes haywire,
it shouldn't affect production applications that are running in the cluster.

In Deis, deployed applications can be segregated by using the ``deis tags`` command. This
enables you to tag machines in your cluster with arbitrary metadata, then configure your applications
to be scheduled to machines which match the metadata.

For example, if some machines in your cluster are tagged with ``environment=production`` and some
with ``environment=staging``, you can configure an application to be deployed to the production
environment by using ``deis tags set environment=production``. Deis will pass this configuration
along to the scheduler, and your applications in different environments on running on separate
hardware.

.. _deis_on_public_clouds:

Running Deis on Public Clouds
-----------------------------
If you are running on a public cloud without security group features, you will have to set up
security groups yourself through either ``iptables`` or a similar tool. The only ports that should
be exposed to the public are:

 - 22: for remote SSH
 - 80: for the routers
 - 443: (optional) routers w/ SSL enabled
 - 2222: for the builder

For providers that do not supply a security group feature, please try
`contrib/util/custom-firewall.sh`_.

.. note::
    If you need to add a new node to the cluster and you are using the custom firewall 
    `contrib/util/custom-firewall.sh`_ you must allow the access to the cluster running
    the next command in each existing node:

.. code-block:: console

    $ NEW_NODE="IP address" contrib/util/custom-firewall.sh

Router firewall
---------------
The :ref:`Router` component includes a firewall to help thwart attacks. It can be enabled by running:
``deisctl config router set firewall/enabled=true``. For more information, see the `router README`_
and :ref:`router_settings`.

IP Whitelist
------------
You can enforce cluster-wide IP whitelisting by running ``deisctl config router set enforceWhitelist=true``.
Then you'll have to manually whitelist IPs to the applications using the config endpoint of the deis
client. The format is ``{IP_or_CIDR}:{Optional_label},...``. For example:

.. code-block:: console

    $ deis config:set -a your-app DEIS_WHITELIST="10.0.1.0/24:office_ABC,212.121.212.121:client_YXZ"

The format is the same for the controller whitelist but you need to specify the list directly into
ectd. For example:

.. code-block:: console

    $ deisctl config router set controller/whitelist="10.0.1.0/24:office_intranet,121.212.121.212:dev_jenkins"


.. _`#986`: https://github.com/deis/deis/issues/986
.. _`contrib/util/custom-firewall.sh`: https://github.com/deis/deis/blob/master/contrib/util/custom-firewall.sh
.. _`router README`: https://github.com/deis/deis/blob/master/router/README.md
