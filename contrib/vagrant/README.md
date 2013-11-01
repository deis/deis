Local development workflow for Deis
================================================================

This document suggests some ways that might make developing Deis easier and cheaper. You don't have to follow
them all.

1. You'll need VirtualBox >= `4.2.18`. We recommend installing Vagrant with their binary installer from http://downloads.vagrantup.com
Vagrant 1.3.5 has support for VirtualBox 4.3

2. Firstly you need to decide whether to use your own Chef Server or the free tier of the hosted enterprise 
service from Opscode. The free tier has a limit of 5 nodes which is more then enough for development. Also
bear in mind that a local Chef Server VM will take up at least 1GB of RAM.

    **Local Chef Server**
    * `cd [DEIS_DIR] && ln -s contrib/vagrant/knife-config .chef`
    * `vagrant up` the chef server Vagrantfile.
    * copy the admin.pem and validation.pem files for your own knife client
    `scp -r root@chefserver.local:/etc/chef-server/admin.pem [DEIS_DIR]/contrib/vagrant/knife-config/`
    `scp -r root@chefserver.local:/etc/chef-server/chef-validator.pem [DEIS_DIR]/contrib/vagrant/knife-config/`

    **Hosted Chef Server**
    * Goto https://getchef.opscode.com/signup and fill in your details.
    * Goto https://preview.opscode.com/login and sign in to your Chef Server.
    * Click on the 'Administration' tab and choose your organisation. There should be a tab in the sidebar that says
    'Starter Kit'. Click it and it will start a small download.
    * Inside the Starter Kit there is a '.chef' folder. Copy it to the root of your Deis codebase.

3. Now you can follow the standard deis setup:
  ```bash
  bundle install # Installs gem files like the knife tool
  berks install # Downloads the relevant cookbooks
  # '--ssl-verify' is only needed when using a local Chef Server
  berks upload [--ssl-verify=false] # Upload the cookbooks to the Chef Server
  ```

4. Use the provision script to boot the deis controller.
    * If you don't already have the deis-node Vagrant box installed (~1GB). This step might take a long time! If for some reason 
    you want to manually add it, use:
    `vagrant box add deis-node https://s3-us-west-2.amazonaws.com/opdemand/deis-node.box`
    * `cd contrib/vagrant && ./provision-controller.sh`
    * You will be asked to add the Controller's SSH key to your local SSH server. This will allow the Controller
    to run vagrant commands on your machine to bootstrap new nodes.
    * If you are using a local Chef Server you will need to tell it that your new controller has permission to create
    nodes. Use:
    `knife client edit deis-controller`
    and your default text editor will launch, you need to set 'admin' to 'true'.

5. The Controller needs to be able to run Vagrant commands on your host machine. It does this via SSH. Therefore
you will need a running SSH server open on port 22.
    * On Debian-flavoured Linux you just need to;
    `sudo apt-get install openssh-server`
    * On Mac OSX you just need to go to **System Preferences -> Sharing** and enable 'Remote Login'.
    * **NB** If your machine's IP changes the Controller won't be able to run commands any more. Currently you aren't informed
    of this, you just get a 500 error from the client. If you tail `/var/log/deis/celeryd.log` though you'll know. When this
    does happen just reupload your IP with `deis providers:discover`.

6. If you want to hack on the actual codebase, you can mount your local codebase onto the VM
   by using the custom Vagrantfile.local.
   * `cp Vagrantfile.local.example Vagrantfile.local` (don't worry it's in .gitignore)
   * Update the VM with `vagrant reload --provision`
   * When mounted you can use your favourite editor to change the code _on your local machine's path_ and then run
   `service deis-server restart` and/or `service deis-worker restart` on the VM for your changes to instantly take effect.
   * It's worth having a read of `Vagrantfile.local.example`

7. If you want to hack on the command line client (`/client/deis.py`), install your local dev version rather than
the one from Pip.
    * `cd deis && make install` This will symlink the dev version to your executables path.
    * Your deis controller is available at http://deis-controller.local so you can register with;
    `deis register http://deis-controller.local`

8. Right, time to boot up some nodes!
  * Create a foramtion with a vagrant flavour, 512MB, 1024MB and 2048MB are available.
  `deis formations:create dev --flavor=vagrant-512 --domain=deisapp.local`
  * Scale a node with `deis nodes:scale dev runtime=1` Be patient, this is the command that runs vagrant commands. Scaling a single node
  can take about 5 mins.
  * Then create and push your app as per the usual documentation.

## Useful development commands
* To use Django's manage.py:
  * SSH in to the VM with `vagrant ssh`
  * Switch user to deis with `sudo su deis`
  * `cd /opt/deis/controller` and activate Venv with `. venv/bin/activate`
  * Get a list of commands with; `./manage.py help`.

* To reset the DB:
  * On the VM run `sudo su postgres -c 'dropdb deis && createdb --encoding=utf8 --template=template0 deis'`
  * When you restart the server with `sudo service deis-server restart` Django will reinstall the DB.
  * You'll need to reupdate the Django's Site Object
  `sudo su deis -c "psql deis -c \"UPDATE django_site SET domain = 'deis-controller.local', name = 'deis-controller.local' WHERE id = 1 \""`

* This is useful for uploading your own local version of the cookbooks, rather than the Github versions:
  * `knife cookbook upload deis --cookbook-path [deis-cookbook path] --force`
  * You need to change directory structure though as knife gets the cookbook name from the folder name. So, for example I use;

  ```bash
  deis
    |__ code
      |__ api
      |__ bin
      |__ client
      |__ cm
      |__ ... and so on
    |__ cookbooks
      |__ deis
        |__ attributes
        |__ definitions
        |__ recipes
        |__ ... and so on
  ```

Notes
-----

Mac OS X: if you see an error such as
'failed to open /dev/vboxnetctl', try restarting VirtualBox:
sudo /Library/StartupItems/VirtualBox/VirtualBox restart