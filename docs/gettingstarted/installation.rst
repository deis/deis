:title: Installation
:description:  Install your own Deis platform on EC2. Configure the Chef server, provision a Deis controller, install the Deis Client and create & scale a formation.

.. _installation:

Installation
============

Follow the steps below to install your own Deis platform on EC2. To complete the
installation process, you will need `Git`_, `RubyGems`_, `Pip`_, the `Amazon EC2 API Tools`_,
`EC2 Credentials`_, and a Chef Server with a working `Knife`_ client.

*Please note: Deis is still under active development. It should not yet be used in production.*

1. Clone the Deis Repository
----------------------------

.. code-block:: console

    $ git clone https://github.com/opdemand/deis.git
    $ cd deis

Cloning the default master branch will provide you with the latest development version
of Deis.  If you want to deploy the latest stable release, make sure you checkout the
most recent tag using ``git checkout vX.Y.Z``.

2. Configure the Chef Server
----------------------------

Deis requires a Chef Server. `Sign up for a free Hosted Chef account`_ if you don’t have one.
You’ll also need a `Ruby`_ runtime with `RubyGems`_ in order to install the required
Ruby dependencies.

.. code-block:: console

    $ bundle install    # install ruby dependencies
    $ berks install     # install cookbooks into your local berkshelf
    $ berks upload      # upload cookbooks to the chef server


3. Provision a Deis Controller
------------------------------

The `Amazon EC2 API Tools`_ will be used to setup basic EC2 infrastructure.  The
`Knife`_ EC2 plugin will be used to bootstrap the controller.

.. code-block:: console

    $ contrib/provision-ec2-controller.sh
    usage: contrib/provision-ec2-controller.sh [region]
    $ contrib/provision-ec2-controller.sh us-west-2
    Creating security group: deis-controller
    + ec2-create-group deis-controller -d 'Created by Deis'
    GROUP    sg-7c40f317    deis-controller    Created by Deis
    + set +x
    Authorizing TCP ports 22,80,443,514 from 0.0.0.0/0...
    + ec2-authorize deis-controller -P tcp -p 22 -s 0.0.0.0/0
    + ec2-authorize deis-controller -P tcp -p 80 -s 0.0.0.0/0
    + ec2-authorize deis-controller -P tcp -p 443 -s 0.0.0.0/0
    + ec2-authorize deis-controller -P tcp -p 514 -s 0.0.0.0/0
    + set +x
    Creating new SSH key: deis-controller
    + ec2-create-keypair deis-controller
    + chmod 600 /home/myuser/.ssh/deis-controller
    + set +x
    Saved to /home/myuser/.ssh/deis-controller
    Created data_bag[deis-users]
    Created data_bag[deis-formations]
    Updated data_bag[deis-apps]
    Provisioning deis-controller with knife ec2...
    ...

Once the ``deis-controller`` node exists on the Chef server, you must log in to
the WebUI add deis-controller to the ``admins`` group. This is required so the
controller can delete node and client records during future scaling operations.


4. Install the Deis Client
--------------------------

Install the Deis client using `Pip`_ (for latest stable) or by linking
``<repo>/client/deis.py`` to ``/usr/local/bin/deis`` (for dev version).
Registration will discover SSH keys automatically and use the 
`standard environment variables`_ to configure the EC2 provider.

.. code-block:: console

    $ sudo pip install deis
    $ deis register http://my-deis-controller.fqdn
    username: myuser
    password:
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done

    Found EC2 credentials: AKIAJTVXXXXXXXXXXXXX
    Import these credentials? (y/n) : y
    Uploading EC2 credentials... done


5. Create & Scale a Formation
-----------------------------

Use the Deis client to create a new formation named "dev" that
has a default layer that serves as both runtime (hosts containers)
and proxy (routes traffic to containers).  Scale the default layer
up to two nodes.

.. code-block:: console

    $ deis formations:create dev --flavor=ec2-us-west-2
    Creating formation... done, created dev
    
    Creating runtime layer... done in 1s
    
    Use `deis nodes:scale dev runtime=1` to scale a basic formation
    
    $ deis nodes:scale dev runtime=2
    Scaling nodes... but first, coffee!
    ...done in 251s
    
    Use `deis create --formation=dev` to create an application


6. Deploy & Scale an Application
--------------------------------

Find an application you’d like to deploy, or clone `an example app`_.
Change into the application directory and use  ``deis create --formation=dev`` 
to create a new application attached to the dev formation.

To deploy the application, new ``git push deis master``.  
Deis will automatically deploy Docker containers 
and configure Nginx proxies to route requests to your application.

Once your application is deployed, you use ``deis scale web=4`` to 
scale up web containers.  You can also use ``deis logs`` to view 
aggregated application logs, or ``deis run`` to run one-off admin
commands inside your application.

To learn more, use ``deis help`` or browse `the documentation`_.

.. code-block:: console

    $ deis create --formation=dev
    Creating application... done, created peachy-waxworks
    Git remote deis added
    $ git push deis master
    Counting objects: 146, done.
    Delta compression using up to 8 threads.
    Compressing objects: 100% (122/122), done.
    Writing objects: 100% (146/146), 21.54 KiB, done.
    Total 146 (delta 84), reused 47 (delta 22)
           Node.js app detected
    -----> Resolving engine versions
           Using Node.js version: 0.10.17
           Using npm version: 1.2.30
    ...
    -----> Building runtime environment
    -----> Discovering process types
           Procfile declares types -> web

    -----> Compiled slug size: 4.7 MB
           Launching... done, v2

    -----> peachy-waxworks deployed to Deis
           http://ec2-54-214-143-104.us-west-2.compute.amazonaws.com ...

    $ curl -s http://ec2-54-214-143-104.us-west-2.compute.amazonaws.com
    Powered by Deis!
    
    $ deis scale web=4
    Scaling containers... but first, coffee!
    done in 12s
    
    === peachy-waxworks Containers
    
    --- web: `node server.js`
    web.1 up 2013-09-23T19:02:30.745Z (dev-runtime-2)
    web.2 up 2013-09-23T19:36:48.741Z (dev-runtime-1)
    web.3 up 2013-09-23T19:36:48.758Z (dev-runtime-1)
    web.4 up 2013-09-23T19:36:48.771Z (dev-runtime-2)


.. _`Git`: http://git-scm.com
.. _`RubyGems`: http://rubygems.org/pages/download
.. _`Pip`: http://www.pip-installer.org/en/latest/installing.html
.. _`Amazon EC2 API Tools`: http://aws.amazon.com/developertools/Amazon-EC2/351
.. _`EC2 Credentials`: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SettingUp_CommandLine.html#set_aws_credentials_linux
.. _`Knife`: http://docs.opscode.com/knife.html
.. _`Sign up for a free Hosted Chef account`: https://getchef.opscode.com/signup
.. _`Ruby`: http://ruby-lang.org/
.. _`standard environment variables`: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SettingUp_CommandLine.html#set_aws_credentials_linux
.. _`an example app`: https://github.com/opdemand/example-nodejs-express
.. _`the documentation`: http://docs.deis.io/
