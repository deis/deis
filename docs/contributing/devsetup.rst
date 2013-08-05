:title: Developer Setup
:description: Setting up your workstation for Deis development
:keywords: deis, documentation, contributing, developer, setup, chef, knife

.. _devsetup:

Developer Setup
===============

Thank you for contributing to Deis! We have tried to make it simple
to work on Deis, but as an open PaaS there are necessarily several
moving pieces and some setup required. We welcome any suggestions for
automating or simplifying this process.


Prerequisites
-------------

We assume you have a modern UNIX-like environment, such as Linux or
Mac OS X. Many of the tools listed below will probably be present
already, but install those that are not:

- git
- GNU Make
- `Python`_ 2.7.x
- `Ruby`_ 1.9 (for knife.rb)
- `Amazon EC2 API Tools`_

To contribute code back to Deis, you must also have a GitHub.com account
in order to create a pull request.


First Things First
------------------

You should complete all the steps outlined in :ref:`installation` and
ensure your CLI and controller are functional before moving on.


Clone the Deis Cookbook
-----------------------

If you want to modify Deis' Chef recipes, you should also clone its
repository:

.. code-block:: console

	$ git clone -q https://github.com/opdemand/deis-cookbook.git


Make a Virtualenv
-----------------

To keep Deis` requirements separate from other development you may do,
it's preferable to create a **virtual environment** for python.

.. code-block:: console

	$ virtualenv venv --prompt='(deis)'
	New python executable in venv/bin/python
	Installing Setuptools.................................done.
	Installing Pip........................................done.
	$ source venv/bin/activate
	(deis)$ pip install -r dev_requirements.txt --use-mirrors
	Downloading/unpacking azure>=0.7.0 (from -r dev_requirements.txt (line 2))
	  Downloading azure-0.7.0.zip (76kB): 76kB downloaded
	  Running setup.py egg_info for package azure
	Downloading/unpacking boto>=2.9.8 (from -r dev_requirements.txt (line 3))
	...
	Successfully installed azure boto ...
	Cleaning up...
	(deis)$

Make sure you install the requirements in the dev_requirements.txt file,
which contains several additions over the runtime requirements.txt file.
Please see the `virtualenv documentation`_ for more details on python virtual
environments.


Modify Code and Test
--------------------

When changing Python code in the Deis project, keep in mind our :ref:`standards`.

** More to come... **


.. _`virtualenv documentation`: http://www.virtualenv.org/en/latest/
.. _`Python`: http://python.org/
.. _`Ruby`: http://ruby-lang.org/
.. _`Amazon EC2 API Tools`: http://aws.amazon.com/developertools/Amazon-EC2/351
.. _`Knife EC2 plugin`: https://github.com/opscode/knife-ec2
