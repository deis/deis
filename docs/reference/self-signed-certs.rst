:title: Creating a Self-Signed SSL Certificate
:description: How to generate a self-signed certificate for securing your application's endpoints

.. _creating_self_signed_ssl:

Creating a Self-Signed SSL Certificate
======================================

When :ref:`using the app ssl <app_ssl>` feature for non-production applications or when
:ref:`installing SSL for the platform <platform_ssl>`, you can avoid the costs associated with the SSL
certificate by using a self-signed SSL certificate. Though the certificate implements full
encryption, visitors to your site will see a browser warning indicating that the certificate should
not be trusted.


Prerequisites
-------------

The openssl library is required to generate your own certificate. Run the following command in your
local environment to see if you already have openssl installed.

.. code-block:: console

    $ which openssl
    /usr/bin/openssl

If the which command does not return a path then you will need to install openssl yourself:

+----------------+---------------------------------+
| If you have... | Install with...                 |
+================+=================================+
| Mac OS X       | Homebrew: brew install openssl  |
+----------------+---------------------------------+
| Windows        | complete package .exe installed |
+----------------+---------------------------------+
| Ubuntu Linux   | apt-get install openssl         |
+----------------+---------------------------------+


Generate Private Key and Certificate Signing Request
----------------------------------------------------

A private key and certificate signing request are required to create an SSL certificate. These can
be generated with a few simple commands. When the openssl req command asks for a “challenge
password”, just press return, leaving the password empty.

.. code-block:: console

    $ openssl genrsa -des3 -passout pass:x -out server.pass.key 2048
    ...
    $ openssl rsa -passin pass:x -in server.pass.key -out server.key
    writing RSA key
    $ rm server.pass.key
    $ openssl req -new -key server.key -out server.csr
    ...
    Country Name (2 letter code) [AU]:US
    State or Province Name (full name) [Some-State]:California
    ...
    A challenge password []:
    ...


Generate SSL Certificate
------------------------

The self-signed SSL certificate is generated from the server.key private key and server.csr files.

.. code-block:: console

    $ openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

The server.crt file is your site certificate suitable for use with
:ref:`Deis's SSL endpoint <app_ssl>` along with the server.key private key.
