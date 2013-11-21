:title: Bare Metal Installation
:description: How to install a Deis controller on existing hardware and create your private PaaS
:keywords: install, installation, deis, controller, setup, bare metal, hardware

.. _static_installation:

Bare Metal Installation
=======================

.. include:: steps1-2.txt

3. Provision a Deis Controller
------------------------------

To create the controller and nodes, the Deis system expects to find Ubuntu
Linux machines running a Docker-compatible kernel. Ubuntu 12.04.3 LTS 64-bit
is recommended.

The `knife`_ command is used to bootstrap the controller. It should
have been installed by ``bundle install``:

.. code-block:: console

    $ knife bootstrap --help
    knife bootstrap FQDN (options)
            --bootstrap-no-proxy [NO_PROXY_URL|NO_PROXY_IP]
                                         Do not proxy locations for the node being bootstrapped
            --bootstrap-proxy PROXY_URL  The proxy server for the node being bootstrapped
            --bootstrap-version VERSION  The version of Chef to install
        -N, --node-name NAME             The Chef node name for your new node
            --server-url URL             Chef Server URL
    ...

Run `knife`_ to create Deis' data bags:

.. code-block:: console

    $ # create data bags
    $ knife data bag create deis-users 2>/dev/null
    Created data_bag[deis-users]
    $ knife data bag create deis-formations 2>/dev/null
    Created data_bag[deis-formations]
    $ knife data bag create deis-apps 2>/dev/null
    Created data_bag[deis-apps]

Run `knife`_ again with appropriate arguments to bootstrap an existing instance
with chef and install Deis' server components. This takes several minutes
to complete:

.. code-block:: console

    $ # bootstrap the controller with knife
    $ knife bootstrap 198.51.100.22 \
    >  --bootstrap-version 11.6.2 \
    >  --ssh-user ubuntu \
    >  --sudo \
    >  --identity-file ~/.ssh/id_rsa \
    >  --node-name deis-controller \
    >  --run-list "recipe[deis::controller]"
    Bootstrapping Chef on 198.51.100.22
    198.51.100.22 --2013-11-20 15:03:46--  https://www.opscode.com/chef/install.sh
    198.51.100.22 HTTP request sent, awaiting response... 200 OK
    198.51.100.22 Length: 6790 (6.6K) [application/x-sh]
    198.51.100.22 Saving to: `STDOUT'
    198.51.100.22
    ...
    198.51.100.22 Chef Client finished, 74 resources updated
    198.51.100.22
    + set +x
    Please ensure that "deis-controller" is added to the Chef "admins" group.

.. include:: steps3-4.txt

5. Register With the Controller
-------------------------------

Registration will discover SSH keys automatically and use environment
variables to configure supported cloud providers with your credentials.

.. code-block:: console

    $ sudo pip install deis
    $ deis register http://deis.example.com
    username: myuser
    password:
    password (confirm):
    email: myuser@example.com
    Registered myuser
    Logged in as myuser

    Found the following SSH public keys:
    1) id_rsa.pub
    Which would you like to use with Deis? 1
    Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done

6. Deploy a Formation and App
-----------------------------

Bare metal nodes have the same Docker-compatible requirements as the
controller does: Ubuntu 12.04.3 LTS 64-bit is recommended.

Create a "static" formation:

.. code-block:: console

    $ deis formations:create dev --flavor=static
    Creating formation... done, created dev
    Creating runtime layer... done in 1s
    $ # if necessary, update runtime layer contents to access your nodes
    $ deis layers:update dev runtime --ssh_username=myuser

Prepare the node for provisioning by the controller:

.. code-block:: console

    $ # use some command-line wizardry to capture just the public key
    $ ssh_key=$(deis layers:info dev runtime | \
    >           grep -Eo '\"ssh_public_key\"\: \"(.*)\"' | \
    >           cut -d\" -f4)
    $ authfile=.ssh/authorized_keys
    $ tmpfile=/tmp/authorized_keys.tmp
    $ # prepend the layer's public key to the node's authorized_keys file
    $ ssh myuser@node1.example.com \
    >  "echo $ssh_key|cat - $authfile > $tmpfile && mv $tmpfile $authfile"

Scale up the formation by adding the existing node:

.. code-block:: console

    $ # add the node to the formation
    $ deis nodes:create dev node1.example.com --layer=runtime
    Creating node for node1.example.com... done in 107s

.. include:: step6.txt
