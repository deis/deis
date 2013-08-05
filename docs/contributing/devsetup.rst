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


Clone the Deis Repositories
---------------------------

.. code-block:: console

	$ git clone -q https://github.com/opdemand/deis.git
	$ git clone -q https://github.com/opdemand/deis-cookbook.git
	$ cd deis


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


Configure the Chef Server
-------------------------

Deis requires a Chef Server. `Sign up for a free Hosted Chef
account`_ if you don’t have one.  You’ll also need a `Ruby`_ runtime with RubyGems
in order to install the required Ruby dependencies.

.. code-block:: console

	$ bundle install    # install ruby dependencies
	$ berks install     # install cookbooks into your local berkshelf
	$ berks upload      # upload cookbooks to the chef server



Provision a Deis Controller
---------------------------

The `Amazon EC2 API Tools`_ will be used to setup basic EC2 infrastructure.
The `Knife EC2 plugin`_ will be used to bootstrap the controller.

.. code-block:: console

	$ contrib/provision-ec2-controller.sh


.. _`virtualenv documentation`: http://www.virtualenv.org/en/latest/
.. _`Python`: http://python.org/
.. _`Ruby`: http://ruby-lang.org/
.. _`Amazon EC2 API Tools`: http://aws.amazon.com/developertools/Amazon-EC2/351
.. _`Sign up for a free Hosted Chef account`: https://getchef.opscode.com/signup
.. _`Knife EC2 plugin`: https://github.com/opscode/knife-ec2
