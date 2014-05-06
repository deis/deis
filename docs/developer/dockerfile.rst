:title: Deploying with Dockerfiles on Deis
:description: A howto on deploying applications using Dockerfiles

Dockerfiles
===========

A Dockerfile automates the steps you would otherwise take manually to create an image.
Deis supports Dockerfiles right out of the box, so you can run your application in your
own custom Docker image.

Deploy using Dockerfiles
------------------------

With Dockerfiles, the stack you deploy your application upon is limitless.
The only requirement is that your application listens on the $PORT environment variable.
This is so `slugrunner`_ can bind your application to an available port on the runtime host.

For example:

.. code-block:: console

    FROM centos:latest
    MAINTAINER OpDemand <info@opdemand.com>
    ENV PORT 8000
    ADD . /app
    WORKDIR /app
    CMD python -m SimpleHTTPServer $PORT

This will serve your application's root directory on a static file server using Docker's
official CentOS image.  Note the server listens on $PORT, which is defaulted to 8000.

.. _`slugrunner`: https://github.com/deis/slugrunner
