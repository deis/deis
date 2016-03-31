:title: Installing Deis on Bare Metal
:description: How to provision a multi-node Deis cluster on Bare Metal

.. _deis_on_bare_metal:

Bare Metal
==========

Deis clusters can be provisioned anywhere `CoreOS`_ can, including on your own hardware.

Please :ref:`get the source <get_the_source>` while following this documentation.

To get CoreOS running on raw hardware, you can boot with `PXE`_ or `iPXE`_ - this will boot a CoreOS
machine running entirely from RAM. Then, you can `install CoreOS to disk`_.

Check System Requirements
-------------------------

Please refer to :ref:`system-requirements` for resource considerations when choosing a
machine size to run Deis.


Generate SSH Key
----------------

.. include:: ../_includes/_generate-ssh-key.rst


Customize user-data
-------------------


Generate a New Discovery URL
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. include:: ../_includes/_generate-discovery-url.rst


SSH Key
^^^^^^^

Add the public key part for the SSH key generated in the first step to the user-data file:

.. code-block:: console

    ssh_authorized_keys:
      - ssh-rsa AAAAB3... deis


Update $private_ipv4
^^^^^^^^^^^^^^^^^^^^

`CoreOS`_ on bare metal doesn't detect the ``$private_ipv4`` reliably. Replace all occurrences in
the user-data with the (private) IP address of the node.


.. include:: ../_includes/_private-network.rst

Add Environment
^^^^^^^^^^^^^^^

Since `CoreOS`_ doesn't detect private and public IP adresses, ``/etc/environment`` file doesn't
get written on boot. Add it to the `write_files` section of the user-data file:

.. code-block:: console

      - path: /etc/environment
        permissions: 0644
        content: |
          COREOS_PUBLIC_IPV4=<your public ip>
          COREOS_PRIVATE_IPV4=<your private ip>


Install CoreOS to disk
----------------------

Assuming you have booted your bare metal server into `CoreOS`_, you can now perform the
installation to disk.

Review disk usage
^^^^^^^^^^^^^^^^^

See :ref:`disk_usage` for more information on how to optimize local disks for Deis.

Provide the config file to the installer
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Save the user-data file to your bare metal machine. The example assumes you transferred the config
to ``/tmp/config``


Start the installation
^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: console

    coreos-install -C stable -c /tmp/config -d /dev/sda -V 899.13.0


This will install the latest `CoreOS`_ stable release that has been known to work
well with Deis. The Deis team tests each new stable release for Deis compatibility,
and it is generally not recommended to use a newer, untested release.

After the installation has finished, reboot your server. Once your machine is back up, you should
be able to log in as the `core` user using the `deis` ssh key.


Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.

Known Problems
--------------


Hostname is localhost
^^^^^^^^^^^^^^^^^^^^^

If your hostname after installation to disk is ``localhost``, set the hostname in user-data before
installation:

.. code-block:: console

    hostname: your-hostname

The hostname must not be the fully qualified domain name!


Slow name resolution
^^^^^^^^^^^^^^^^^^^^

Certain DNS servers and firewalls have problems with glibc sending out requests for IPv4 and IPv6
addresses in parallel. The solution is to set the option ``single-request`` in
``/etc/resolv.conf``. This can best be accomplished in the user-data when installing `CoreOS`_ to
disk. Add the following block to the ``write_files`` section:

.. code-block:: console

      - path: /etc/resolv.conf
        permissions: 0644
        content: |
          nameserver 8.8.8.8
          nameserver 8.8.4.4
          domain your.domain.name
          options single-request


.. _`cluster size`: https://github.com/coreos/etcd/blob/master/Documentation/optimal-cluster-size.md
.. _`CoreOS`: https://coreos.com/
.. _`install CoreOS to disk`: https://coreos.com/docs/running-coreos/bare-metal/installing-to-disk/
.. _`iPXE`: https://coreos.com/docs/running-coreos/bare-metal/booting-with-ipxe/
.. _`PXE`: https://coreos.com/docs/running-coreos/bare-metal/booting-with-pxe/
