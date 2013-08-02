:title: Welcome to the Deis documentation
:description: An overview of the Deis documentation
:keywords: api, commandline, command-line, contributing, faq, terms, tutorial, deis, docker, heroku

Welcome
=======

.. image:: _static/img/deis-graphic.png

Deis is an open source PaaS that makes it easy to deploy
:ref:`container`\s and  :ref:`node`\s used to host applications,
databases, middleware and other services. Deis leverages Chef, Docker and
Heroku Buildpacks to provide a private PaaS that is lightweight and flexible.

If you are new to Deis, start with the :ref:`tutorial`. The :ref:`tutorial`
also references the :ref:`installation` guide.

Once you have set up a Deis Controller, you will use the
``deis`` **command-line client** to create an app, push your code, and
scale your own cloud resources. ``deis --help`` is very helpful, but we
also have a :ref:`cheatsheet` with examples of usage.

If you find a bug, please report it to the project at 
https://github.com/opdemand/deis/issues

Developers who want to explore Deis internals should browse the :ref:`api`. 
Fork the open source deis repository and enjoy the freedom of your own
private PaaS. The ``deis`` project won't ignore pull requests; please help
us improve.
