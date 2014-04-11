:title: Provision a Deis PaaS Controller
:description: Guide to provisioning a Deis controller, the brains of a Deis private PaaS.


.. _provision-controller:

Provision a Controller
======================
The `controller` is the brains of a Deis platform.
There are two ways to provision a controller: automatically or manually.

Automatic Provisioning
----------------------
The community maintains shell scripts that automate the provisioning
of Deis controllers on different cloud providers.
In addition to launching the controller itself, these scripts also
use optimized base images,
generate SSH keys, firewall configs and other cloud infrastructure
per Deis best practices.

You can find instructions on automatic provisioning for:

 * `EC2`_
 * `Rackspace`_
 * `Digital Ocean`_

Please note that even with automatic provisioning, you will still have to
`add the controller to the admins group`_.

Manual Provisioning
-------------------
If you want your controller on bare metal, a different cloud provider,
or would just rather provision things manually --no problem!
Just remember with manual provisioning, you are in charge of:

 * Ensuring system requirements are met
 * SSH key generation and distribution
 * Network configuration

.. important:: System Requirements
   Most controllers require at least 2GB of system memory and 100GB of storage

The general process for manual provisioning involves:

 #. Boot a target host that meets system requirements
 #. Make sure the target host is accessible over SSH from your workstation
 #. Use ``knife bootstrap`` to provision the controller on the target host

Here is an example ``knife bootstrap`` command:

.. code-block:: console

    $ knife bootstrap 198.51.100.22 \
    >  --bootstrap-version 11.8.2 \
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

Please note the ``knife bootstrap`` command can take several minutes to complete.

Add Controller to Admins Group
------------------------------
Whether you used automatic or manual provisioning,
you must add "deis-controller" to the "admins" group on the Chef Server.

Hosted Chef
~~~~~~~~~~~
For Hosted Chef, visit the `OpsCode management interface`_ to add
"deis-controller" to the "admins" group. These `steps`_ describe specifically
how to add users to groups in the management interface.

Open Source Chef
~~~~~~~~~~~~~~~~

For Open Source Chef, use the knife cli tool. First, set the EDITOR environment
variable to your favourite text editor:

``$ export EDITOR=vim``

Now issue the command to edit the "deis-controller" user.

``$ knife client edit deis-controller``

and update the 'admin' key to 'true'. The updated content should look similar to this.

.. code-block:: json

  {
    "name": "deis-controller",
    "public_key": "-----BEGIN PUBLIC KEY-----\nMaa...QAB\n-----END PUBLIC KEY-----\n",
    "admin": true,
    "json_class": "Chef::ApiClient",
    "chef_type": "client"
  }

Save and close the file. The "deis-controller" user is now in the "admins" group.

.. important::
   If you skip adding "deis-controller" to the Chef "admins" group, you will
   receive errors when scaling down nodes, as the controller will not have
   permission to delete "client" and "node" records from the Chef Server.

.. _`EC2`: https://github.com/opdemand/deis/tree/master/contrib/ec2#readme
.. _`Rackspace`: https://github.com/opdemand/deis/tree/master/contrib/rackspace#readme
.. _`Digital Ocean`: https://github.com/opdemand/deis/tree/master/contrib/digitalocean#readme
.. _`add the controller to the admins group`: #add-controller-to-admins-group
.. _`knife`: http://docs.opscode.com/knife.html
.. _`OpsCode management interface`: https://manage.opscode.com/
.. _`steps`: http://docs.opscode.com/manage_server_hosted_groups.html#add-user-to-group

