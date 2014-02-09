Standalone Service Registry
===========================

[![Build Status](https://travis-ci.org/bacongobbler/service-registry.png?branch=master)](https://travis-ci.org/bacongobbler/service-registry)

The service registry is a wrapper to common web services for
centralized service creation (known as provisioning). It can work with
common on-premise services such as a local MySQL database or a Redis server
for a "private" or on-premise service registry, or with public SaaS applications
such as MongoLab or ElephantSQL for a public service offering.

## Technical Overview

This project will be a generic wrapper for web services, giving administrators
the power to make their consumable service readily available to any PaaS that
uses this service registry.

Similar to Heroku Add-ons, customers use the service registry to provision an
service offered by a third-party provider for their applications. When this
happens, the service registry sends a request to a providers's service provider,
which then creates a new private service. This service represents an instance of
the provider's service.

## Target Market

The service registry aims to target the PaaS industry, where new PaaS vendors
are looking to add data services to their PaaS offering (similar to Heroku
Add-Ons). Private PaaS solutions may also implement this solution by
connecting to external databases or services instead, whether those services
are internal to a clustered environment or are in a large datacenter in another
region.

It also targets the SaaS market, where new providers are looking to get some
public visibility of their product through third parties that implement
this service in their offering.

## Registry Endpoints

    GET     /services
        retrieve information about all providers available

    POST    /services
        attach a provider to the registry

    GET     /services/:provider_name
        get information about this specific provider

    PUT     /services/:provider_name
        update configuration for a provider

    DELETE  /services/:provider_name
        unregister a provider. also unprovisions all services provisioned with this provider

    GET     /services
        get information on all services available

    GET     /services?user=:username
        get all services provisioned by :username

    GET     /services?provider=:provider_name
        get all services provisioned by a specific provider

    POST    /services
        provisions a new service

    GET     /services/:service_name
        get information about the service

    PUT     /services/:service_name
        update an service's config

    DELETE /services/:service_name
        deprovision an service

## Architecture Diagram

    +-----------------+
    |                 |
    | External Client |-----+
    |                 |     |
    +-----------------+     |   POST /services { type: mysql, plan: free }
                            |   RETURN { status: 200, creds: { MYSQL_URL: ... }, name: mysql-746017aa }
    192.168.0.1             |                   |
    +------------+          | <-----------------+
    |            |          |
    |  Registry  |----------+
    |            |         CREATE DATABASE 746017aa; GRANT USER ...;
    +------------+                                   |
          |                                          |
          | POST / { plan: free }                    |
          | RETURN { status: 200, creds: {...} }     |
          |            |                             |
          |            |           192.168.0.2       |      192.168.0.3:3306
          |            |        +---------------+    |      +---------------+
          |            V        |               |    V      |               |
          +---------------------|   Provider    +-----------+   MySQL       |
          |                     |               |           |               |
          |                     +---------------+           +---------------+
          |                        192.168.0.4              192.168.0.3:5432
          |                     +---------------+           +---------------+
          |                     |               |           |               |
          +---------------------|   Provider    +-----------+   PostgreSQL  |
          |                     |               |           |               |
          |                     +---------------+           +---------------+
          |                        192.168.0.5              http://mongolab.com/
          |                     +---------------+           +---------------+
          |                     |               |           |               |
          +---------------------|   Provider    +-----------+   MongoLab    |
                                |               |           |               |
                                +---------------+           +---------------+

# Testing/Contributing

You can test the project by running:

    $ pip install tox
    $ tox

If you want to contribute to this project, please feel free to
[make a pull request](https://github.com/bacongobbler/service-registry/pulls) or
[open an issue](https://github.com/bacongobbler/service-registry/issues). If you're
making a pull request, please supply tests and make sure that `tox` passes without
failure.
