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
    * Goto https://manage.opscode.com/login and sign in to your Chef Server.
    * Click on the 'Administration' tab and choose your organisation. There should be a tab in the sidebar that says
    'Starter Kit'. Click it and it will start a small download.
    * Inside the Starter Kit there is a '.chef' folder. Copy it to the root of your Deis codebase.

3. Now you can follow the standard deis setup:
  * If you're running a local chef server, you should adjust the `Gemfile` and make sure the version of berkshelf is 3.0.x. This is needed for the `--ssl-verify` option to work correctly.
  ```bash
  bundle install # Installs gem files like the knife tool
  berks install # Downloads the relevant cookbooks
  # '--ssl-verify' is only needed when using a local Chef Server
  berks upload [--ssl-verify=false] # Upload the cookbooks to the Chef Server
  ```

4. The Controller needs to be able to run Vagrant commands on your host machine. It does this via SSH. Therefore
you will need a running SSH server open on port 22 and a means to broadcast your hostname to local DNS.
    * On Debian-flavoured Linux you just need to;
    `sudo apt-get install openssh-server`
    * On Mac OSX you just need to go to **System Preferences -> Sharing** and enable 'Remote Login'.
    * [Mac OSX user's should already be broadcasting their hostname](http://support.apple.com/kb/ht3473).
    * Linux user's will need to install avahi-daemon, so that their machine is accessible via
    [hostname].local. Eg; `sudo apt-get install avahi-daemon`.

5. Creating the Deis Controller.
    * Run `./contrib/vagrant/provision-vagrant-controller.sh`
    * You may need to prepend the command with `bundle exec` depending on your Ruby setup.
    * When running for the first time you will be asked to add the Controller's SSH key to your local SSH server.
    This will allow the Controller to run vagrant commands on your machine to bootstrap new nodes.
    * You need to tell the Chef Server that your new Controller has permission to create
    and delete nodes. Use:
      * For a local Chef Server just type `knife client edit deis-controller` and your default text
      editor will launch, you need to set 'admin' to 'true'.
      * For Hosted Chef, log in to https://manage.opscode.com/. Go to the
      Administration tab, click on the "Groups" entry to the left, then the "admins" entry
      under "All Groups". Choose the "Permissions" tab and click the "+ Add" button, then
      type in "deis-controller" and add it. Assign all permissions to the "deis-controller"
      client object.

6. If you want to hack on the command line client (`client/deis.py`), install your local dev version rather than
the one from Pip.
    * `cd deis && make install` This installs the client into your executables path.
    * `sudo rm /usr/local/bin/deis && sudo ln ./deis.py /usr/local/bin/deis` This will symlink the dev version to your executables path.
    * Your deis controller is available at http://deis-controller.local:8000 so you can register with;
    `deis register http://deis-controller.local:8000`

7. Right, time to boot up some nodes!
  * Create a formation with a vagrant flavour, 512MB, 1024MB and 2048MB are available.
  `deis formations:create dev --flavor=vagrant-512 --domain=deisapp.local`
  * Scale a node with `deis nodes:scale dev runtime=1` Be patient, this is the command that runs vagrant commands. Scaling a single node
  can take over 15 mins.
  * Then create and push your app as per the usual documentation.

## Useful development commands
* To get a shell session to a running container use the `dsh` command on the VM. Usage:
  * `dsh deis-builder`
  * `dsh deis-builder /bin/ls` Note the absolute path
  * `echo 'ls' | dsh deis-builder` Note no need for path when piping

* To use Django's manage.py:
  * SSH in to the VM with `vagrant ssh`
  * Use `dsh deis-controller`
  * Get into the Django server's path `cd /app/deis/controller`
  * Get a list of commands with; `./manage.py help`

* Django's native web admin interface is available at http://deis-controller.local:8000/admin/
You can add and update all of the models from there.

* To reset the DB:
  * There is a script at contrib/vagrant/util/reset-db.sh that resets the DB and installs some basic fixtures.
  It should be run from your host machine.
  * It installs a formation named 'dev' and a super user with username 'devuser' and password 'devpass'.

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

* Resolving submodule conflicts when merging master into your dev branch can be troublesome.
You might have luck with `git submodule foreach git checkout master && git submodule foreach git pull --rebase`

* There is a sample `Vagrantfile.local.example` that can be useful if you want any personal customisations
to your vagrant setup, like using less RAM for post-installation boots.

Notes
-----

Mac OS X: if you see an error such as
'failed to open /dev/vboxnetctl', try restarting VirtualBox:
sudo /Library/StartupItems/VirtualBox/VirtualBox restart
