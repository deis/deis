:title: Deploying with Dockerfiles on Deis
:description: A howto on deploying applications using Dockerfiles
:keywords: tutorial, guide, walkthrough, howto, deis, developer, dev, docker, dockerfile

Dockerfiles
===========

A Dockerfile automates the steps you would otherwise take manually to create an image.
Deis supports Dockerfiles right out of the box, so you can run your application in your
own custom Docker image.

Deploy using Dockerfiles
------------------------

With Dockerfiles, the stack you deploy your application upon is limitless. The only
requirement is that it sets a ENV entry to set the PORT environment variable. This is so
`slugrunner`_ can listen for when the process is alive or dead. For example:

.. code-block:: console

    FROM centos:latest
    MAINTAINER OpDemand <info@opdemand.com>
    ENV PORT 8000
    ADD . /app
    WORKDIR /app
    CMD python -m SimpleHTTPServer $PORT

Which will serve your application's root directory on a static file server on Docker's
official CentOS image.

.. _`slugrunner`: https://github.com/deis/slugrunner
