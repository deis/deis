:title: Using a Proxy Server
:description: How to configure Deis to use a proxy server.

.. _using-a-proxy-server:

Using a Proxy Server
====================

In some environments, HTTP connections must pass through a proxy. The Deis builder component supports
proxies by respecting the proxy-related environment variables defined in ``/etc/environment_proxy``.

Additionally, Docker is also configured to respect the settings in this file.

By default, ``/etc/environment_proxy`` has all environment variables set to blank values:

.. code-block:: console

    HTTP_PROXY=
    HTTPS_PROXY=
    ALL_PROXY=
    NO_PROXY=
    http_proxy=
    https_proxy=
    all_proxy=
    no_proxy=

.. note::

    Proxy settings must be respected by the applications you're building.
    When using custom buildpacks, make sure they respect proxy settings.


Configuring before server launch
--------------------------------

Before provisioning the servers using the provision scripts in the Deis repository, edit
``contrib/coreos/user-data.example`` and replace the contents of the file to suit your environment.

For example:

.. code-block:: console

  - path: /etc/environment_proxy
    owner: core
    content: |
      HTTP_PROXY=http://proxy.example.com:3128
      HTTPS_PROXY=http://proxy.example.com:3128
      ALL_PROXY=http://proxy.example.com:3128
      NO_PROXY="127.0.0.1,localhost,.example.com"
      http_proxy=http://proxy.example.com:3128
      https_proxy=http://proxy.example.com:3128
      all_proxy=http://proxy.example.com:3128
      no_proxy="127.0.0.1,localhost,.example.com"

After running ``make discovery-url`` and provisioning your servers, the platform will come up with
your proxy settings.

Configuring after server launch
-------------------------------

It's also possible to configure these settings after the server has been provisioned, but this will
result in downtime of the Deis platform as components are restarted.

You'll need to edit ``/etc/environment_proxy`` on all CoreOS hosts (as the builder component can
be relocated to any host in the cluster). Then, restart Docker with ``sudo systemctl restart docker``
and monitor Deis components with ``deisctl list``. It may be necessary to restart components
if they do not recover automatically from the Docker restart.
