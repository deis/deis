:title: Using an SSL Certificate with Deis
:description: Enabling and configuring SSL on applications using the SSL endpoint.


.. _app_ssl:

Application SSL Certificates
============================

SSL is a cryptographic protocol that provides end-to-end encryption and integrity for all web
requests. Apps that transmit sensitive data should enable SSL to ensure all information is
transmitted securely.

To enable SSL on a custom domain, e.g., ``www.example.com``, use the SSL endpoint.

.. note::

    ``deis certs`` is only useful for custom domains. Default application domains are
    SSL-enabled already and can be accessed simply by using https,
    e.g. ``https://foo.deisapp.com`` (provided that you have :ref:`installed your wildcard
    certificate <router_ssl>` on the routers or :ref:`on the load balancer <load_balancer_ssl>`).


Overview
--------

Because of the unique nature of SSL validation, provisioning SSL for your domain is a multi-step
process that involves several third-parties. You will need to:

1. Purchase an SSL certificate from your SSL provider
2. Upload the cert to Deis


Acquire SSL Certificate
-----------------------

Purchasing an SSL cert varies in cost and process depending on the vendor. `RapidSSL`_ offers a
simple way to purchase a certificate and is a recommended solution. If you’re able to use this
provider, see `buy an SSL certificate with RapidSSL`_ for instructions.


DNS and Domain Configuration
----------------------------

Once the SSL certificate is provisioned and your cert is confirmed, you must route requests for
your domain through Deis. Unless you've already done so, add the domain specified when generating
the CSR to your app with:

.. code-block:: console

    $ deis domains:add www.example.com -a foo
    Adding www.example.com to foo... done


Attach the Certificate
----------------------

Add your certificate, any intermediate certificates, and private key to the endpoint with the
``certs:add`` command.

.. code-block:: console

    $ deis certs:add server.crt server.key
    Adding SSL endpoint... done
    www.example.com

.. note::

    It may take up to one minute for the certificate to be available on the routers.


Attach a Certificate Chain
^^^^^^^^^^^^^^^^^^^^^^^^^^

Sometimes, your certificates (such as a self-signed or a cheap certificate) need additional
certificates to establish the chain of trust. What you need to do is bundle all the certificates
into one file and give that to Deis. Importantly, your site’s certificate must be the first one:

.. code-block:: console

    $ cat server.crt server.ca > server.bundle

After that, you can add them to Deis with the ``certs:add`` command:

.. code-block:: console

    $ deis certs:add server.bundle server.key
    Adding SSL endpoint... done
    www.example.com


Endpoint Details
----------------

You can verify the details of your domain's SSL configuration with ``deis certs``.

.. code-block:: console

    $ deis certs
    Common Name      Expires
    ---------------  ----------------------
    www.example.com  2016-12-31T00:00:00UTC


Testing SSL
-----------

Use a command line utility like ``curl`` to test that everything is configured correctly for your
secure domain.

.. note::

    The -k option flag tells curl to ignore untrusted certificates.

Pay attention to the output. It should print ``SSL certificate verify ok``. If it prints something
like ``common name: www.example.com (does not match 'www.somedomain.com')`` then something is not
configured correctly.

Remove Certificate
------------------

You can remove a certificate using the ``certs:remove`` command:

.. code-block:: console

    $ deis certs:remove www.example.com
    Removing www.example.com... Done.


Troubleshooting
---------------

Here are some steps you can follow if your SSL endpoint is not working as you'd expect.


Untrusted Certificate
^^^^^^^^^^^^^^^^^^^^^

In some cases when accessing the SSL endpoint, it may list your certificate as untrusted.

If this occurs, it may be because it is not trusted by Mozilla’s list of `root CAs`_. If this is
the case, your certificate may be considered untrusted for many browsers.

If you have uploaded a certificate that was signed by a root authority but you get the message that
it is not trusted, then something is wrong with the certificate. For example, it may be missing
`intermediary certificates`_. If so, download the intermediary certificates from your SSL provider,
remove the certificate from Deis and re-run the ``certs:add`` command.

.. _`RapidSSL`: https://www.rapidssl.com/
.. _`buy an SSL certificate with RapidSSL`: https://www.rapidssl.com/buy-ssl/
.. _`root CAs`: https://www.mozilla.org/en-US/about/governance/policies/security-group/certs/included/
.. _`intermediary certificates`: http://en.wikipedia.org/wiki/Intermediate_certificate_authorities
