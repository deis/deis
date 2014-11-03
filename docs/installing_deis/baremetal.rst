:title: Installing Deis on Bare Metal
:description: How to provision a multi-node Deis cluster on Bare Metal

.. _deis_on_bare_metal:

Bare Metal
==========

Deis clusters can be provisioned anywhere `CoreOS`_ can, including on your own hardware. To get
CoreOS running on raw hardware, you can boot with `PXE`_ or `iPXE`_ - this will boot a CoreOS
machine running entirely from RAM. Then, you can `install CoreOS to disk`_.


Generate SSH key
----------------

To avoid problems deploying/launching apps later on it is necessary to install `CoreOS`_ to disk
with a SSH key without a passphrase. The following command will generate a new keypair named
"deis":

.. code-block:: console

    $ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis


Customize user-data
-------------------


Generate a New Discovery URL
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

To get started with provisioning Deis, we will need to generate a new Discovery URL. Discovery URLs
help connect `etcd`_ instances together by storing a list of peer addresses and metadata under a
unique address. You can generate a new discovery URL for use in your platform by
running the following from the root of the repository:

.. code-block:: console

    $ make discovery-url

This will write a new discovery URL to the user-data file. Some convenience scripts are supplied in
this user-data file, so it is mandatory for provisioning Deis.


SSH Key
^^^^^^^

Add the public key part for the SSH key generated in the first step to the user-data file:

.. code-block:: console

    ssh_authorized_keys:
      - ssh-rsa AAAAB3... deis


Update $private_ipv4
^^^^^^^^^^^^^^^^^^^^

`CoreOS`_ on bare metal doesn't detect the ``$private_ipv4`` reliably. Replace all occurences in
the user-data with the (private) IP address of the node.


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

Assuming you have booted your bare metal server into `CoreOS`_, you can perform now perform the
installation to disk.

Provide the config file to the installer
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Save the user-data file to your bare metal machine. The example assumes you transferred the config
to ``/tmp/config``


Start the installation
^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: console

    coreos-install -C alpha -c /tmp/config -d /dev/sda


This will install the current `CoreOS`_ release to disk. If you want to install the recommended
`CoreOS`_ version, check the `Deis changelog`_ and specify that version by appending the ``-V``
parameter to the install command, e.g. ``-V 472.0.0``.

After the installation has finished, reboot your server. Once your machine is back up, you should
be able to log in as the `core` user using the `deis` ssh key.


Configure DNS
-------------

See :ref:`configure-dns` for more information on properly setting up your DNS records with Deis.


Install Deis Platform
---------------------

Now that you've finished provisioning a cluster, please refer to :ref:`install_deis_platform` to
start installing the platform.


Considerations when deploying Deis
----------------------------------

* Use machines with ample disk space and RAM (for comparison, we use m3.large instances on EC2)
* Choose an appropriate `cluster size`_


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
.. _`Deis changelog`: https://github.com/deis/deis/blob/master/CHANGELOG.md
.. _`etcd`: https://github.com/coreos/etcd
.. _`install CoreOS to disk`: https://coreos.com/docs/running-coreos/bare-metal/installing-to-disk/
.. _`iPXE`: https://coreos.com/docs/running-coreos/bare-metal/booting-with-ipxe/
.. _`PXE`: https://coreos.com/docs/running-coreos/bare-metal/booting-with-pxe/