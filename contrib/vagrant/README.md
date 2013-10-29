Provision a Deis Controller on Vagrant
======================================

This document describes how to set up a Deis controller and two nodes with
Vagrant for testing.

1. Install VirtualBox version 4.2.18. (Vagrant does not support version 4.3.)
Then start VirtualBox and install the VirtualBox Extension Pack for 4.2.18.

2. Install Vagrant version 1.3.5 or later. Rather than fight Ruby dependencies,
use a binary installer from vagrantup.com.

3. Run the provisioning script:
```console
$ ./contrib/vagrant/provision-vagrant-controller.sh
```

This script will:
- Create the data bags in your Chef account to support Deis
- Run `vagrant up` to create a Deis controller and 2 static nodes
- Register the controller with Chef and install Deis and supporting software

4. Register a user

5. Run the script to create a static formation:
```console
$ ./contrib/vagrant/create-static-formation.sh
```

Notes
-----

Mac OS X: if you see an error such as
"failed to open /dev/vboxnetctl", try restarting VirtualBox:
sudo /Library/StartupItems/VirtualBox/VirtualBox restart
