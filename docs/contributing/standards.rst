:title: Coding Standards
:description: -- todo: change me
:keywords: deis, documentation, contributing, developer, setup

.. _standards:

Coding Standards
================

Deis is a `python`_ project. We chose python over other compelling
languages because it is widespread, well-documented, and friendly to
a large number of developers. We think open source code benefits from
many eyes upon it.

Contributors to deis should feel welcome to make changes to any part
of the codebase. To create a proper github pull request for inclusion
into the official repository, your code must pass two tests:

- :ref:`flake8`
- :ref:`coverage`


.. _flake8:

``make flake8``
---------------

::

	$ cd $HOME/projects/deis
	$ make flake8
	flake8
	$


.. _coverage:

``make coverage``
-----------------

::

	$ cd $HOME/projects/deis
	$ make coverage
	coverage run manage.py test api celerytasks client web
	Creating test database for alias 'default'...
	...................ss.
	----------------------------------------------------------------------
	Ran 22 tests in 22.630s

	OK (skipped=2)
	Destroying test database for alias 'default'...
	coverage html
	$


.. _python: http://www.python.org/
