:title: Coding Standards
:description: Deis project coding standards. Contributors to Deis should feel welcome to make changes to any part of the codebase.
:keywords: deis, contributing, coding, python, test, tests, testing

.. _standards:

Coding Standards
================

Deis is a `Python`_ project. We chose Python over other compelling
languages because it is widespread, well-documented, and friendly to
a large number of developers. Source code benefits from many eyes
upon it.

`The Zen of Python`_ emphasizes simple over clever, and we agree.
Readability counts. Deis also aims for complete test coverage.

Contributors to Deis should feel welcome to make changes to any part
of the codebase. To create a proper GitHub pull request for inclusion
into the official repository, your code must pass two tests:

- :ref:`make_flake8`
- :ref:`make_coverage`


.. _make_flake8:

``make flake8``
---------------

`flake8`_ is a helpful command-line tool that combines the output of
`pep8 <pep8_tool_>`_, `pyflakes`_, and `mccabe`_.

.. code-block:: console

	$ cd $HOME/projects/deis
	$ make flake8
	flake8
	$

No output, as above, means ``flake8`` found no errors. If errors
are reported, fix them in your source code and try ``flake8`` again.

The Deis project adheres to `PEP8`_, the python code style guide,
with the exception that we allow lines up to 99 characters in length.
Docstrings and tests are also required for all public methods, although
``flake8`` does not enforce this.

Default settings for ``flake8`` are in the ``[flake8]`` section of the
setup.cfg file in the project root.


.. _make_coverage:

``make coverage``
-----------------

Once your code passes the style checker, run the test suite and
ensure that everything passes and that code coverage has not declined.

.. code-block:: console

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

If a test fails, fixing it is obviously the first priority. And if you
have introduced new code, it must be accompanied by unit tests.

In the example above, all tests passed and ``coverage`` created a report
of what code was exercised while the tests were running. Open the file
``htmlcov/index.html`` under the project's root and ensure that the
overall coverage percentage has not receded as a result of your
changes. Current test coverage can be found here:

.. image:: https://coveralls.io/repos/opdemand/deis/badge.png?branch=master
    :target: https://coveralls.io/r/opdemand/deis?branch=master
    :alt: Coverage Status

Now create a GitHub pull request with a description of what your code
fixes or improves. That's it!


.. _Python: http://www.python.org/
.. _flake8: https://pypi.python.org/pypi/flake8/
.. _pep8_tool: https://pypi.python.org/pypi/pep8/
.. _pyflakes: https://pypi.python.org/pypi/pyflakes/
.. _mccabe: https://pypi.python.org/pypi/mccabe/
.. _PEP8: http://www.python.org/dev/peps/pep-0008/
.. _`The Zen of Python`: http://www.python.org/dev/peps/pep-0020/
