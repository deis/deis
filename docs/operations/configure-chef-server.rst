:title: Configure a Chef Server for Deis
:description: Guide to setting up a Chef Server for use with Deis
:keywords: tutorial, guide, walkthrough, howto, deis, sysadmins, operations, chef, chef server

.. _configure-chef-server:

Configure a Chef Server
=======================
The first step in installing Deis is configuring a `Chef Server`_.
Deis relies heavily on Chef Server to provide a battle-tested foundation for:

 * Hosting `cookbooks`_ used to configure Deis :ref:`Nodes <node>`
 * Hosting `data bags`_ used to store global cluster configuration

`Sign up for a free Hosted Chef account`_ if you don't have a Chef Server.

Clone the Deis Repository
-------------------------
The specific cookbook versions required for each Deis release reside in the project's
``Berksfile``.  Use `git`_ to checkout the source code.

.. code-block:: console

    $ git clone https://github.com/opdemand/deis.git
    $ cd deis

Cloning the default master branch will provide you with the latest development
version of Deis. If you instead want to deploy the latest stable release,
checkout the most recent tag using ``git checkout v0.3.0``, for example.

Prepare your Ruby Environment
-----------------------------
You'll need a basic Ruby environment that allows the installation of
`RubyGems`_ like ``knife`` and ``berkshelf``.  
From the Deis checkout directory:

.. code-block:: console

    $ gem install bundler  # install the bundler tool (if necessary)
    $ bundle install       # install ruby dependencies from Gemfile

Test Connectivity to Chef Server
--------------------------------
In order to communicate with the Chef Server, your `knife.rb`_ must be configured properly.
You can test this with:

.. code-block:: console

    $ knife client list
    gabrtv-validator

You should see at least the validator key for your Chef organization.
If not, your `knife.rb`_ configuration or Chef keys are probably incorrect.

Upload Cookbooks
----------------
Deis uses the wonderful `Berkshelf`_ utility to manage cookbooks.
The ``berks`` command should have been installed during ``bundle install``.
To update Deis cookbooks using Berkshelf:

.. code-block:: console

    $ berks install        # install cookbooks to your local berkshelf
    $ berks upload         # upload berkshelf cookbooks to the chef server

Create Data Bags
----------------
Deis uses three data bags on the Chef Server for storing global state
used during chef-client runs:

 * deis-users
 * deis-formations
 * deis-apps

Use `knife`_ to create these data bags on the Chef Server:

.. code-block:: console

    $ knife data bag create deis-users
    Created data_bag[deis-users]
    $ knife data bag create deis-formations
    Created data_bag[deis-formations]
    $ knife data bag create deis-apps
    Created data_bag[deis-apps]

.. _`Chef Server`: http://docs.opscode.com/chef_overview_server.html
.. _`Sign up for a free Hosted Chef account`: https://getchef.opscode.com/signup
.. _`cookbooks`: http://docs.opscode.com/essentials_cookbooks.html
.. _`data bags`: http://docs.opscode.com/essentials_data_bags.html
.. _`knife.rb`: http://docs.opscode.com/config_rb_knife.html
.. _`git`: http://git-scm.com
.. _`RubyGems`: http://rubygems.org/pages/download
.. _`knife`: http://docs.opscode.com/knife.html
.. _`Berkshelf`: http://berkshelf.com/
