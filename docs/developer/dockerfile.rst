:title: Deploying with Dockerfiles on Deis
:description: A howto on deploying applications using Dockerfiles

Dockerfiles
===========

A Dockerfile automates the steps you would otherwise take manually to create an image. Deis
supports Dockerfiles right out of the box, so you can run your application in your own custom
Docker image.

Deploy using Dockerfiles
------------------------

With Dockerfiles, the stack you deploy your application upon is limitless. The only requirements
are:

1) Your application has a Procfile
2) Your Procfile has a web process
3) Your application listens to $PORT
4) Your Dockerfile exposes port 5000.

For example:

.. code-block:: console

    $ cat Dockerfile
    FROM centos:latest
    MAINTAINER OpDemand <info@opdemand.com>
    ADD . /app
    WORKDIR /app
    EXPOSE 5000
    $ cat Procfile
    web: python -m SimpleHTTPServer $PORT


This will serve your application's root directory on a static file server using Docker's official
CentOS image. Note that the web process listens on $PORT.
