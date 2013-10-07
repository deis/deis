Provision a Deis Controller on Rackspace
========================================

1. Install [knife-rackspace][kniferack] with `gem install knife-rackspace` or
just `bundle install` from the root directory of your deis repository:
```console
$ cd $HOME/projects/deis
$ gem install knife-rackspace
Fetching: knife-rackspace-0.8.1.gem (100%)
Successfully installed knife-rackspace-0.8.1
1 gem installed
Installing ri documentation for knife-rackspace-0.8.1...
Installing RDoc documentation for knife-rackspace-0.8.1...
```

2. Export your Rackspace credentials as environment variables and edit knife.rb
to read them:
```console
$ cat <<'EOF' >> $HOME/.bash_profile
export RACKSPACE_USERNAME=<your_rackspace_username>
export RACKSPACE_API_KEY=<your_rackspace_api_key>
EOF
$ source $HOME/.bash_profile
$ cat <<'EOF' >> $HOME/.chef/knife.rb
knife[:rackspace_api_username] = "#{ENV['RACKSPACE_USERNAME']}"
knife[:rackspace_api_key] = "#{ENV['RACKSPACE_API_KEY']}"
EOF
$ knife rackspace server list
Instance ID  Name  Public IP  Private IP  Flavor  Image  State
```

3. Run the provisioning script to create a new Deis controller:
```console
$ ./contrib/rackspace/provision-rackspace-controller.sh dfw
Provisioning a deis controller on Rackspace...
Creating new SSH key: deis-controller
+ ssh-keygen -f /home/myuser/.ssh/deis-controller -t rsa -N '' -C deis-controller
+ set +x
Saved to /home/myuser/.ssh/deis-controller
Created data_bag[deis-users]
Created data_bag[deis-formations]
Created data_bag[deis-apps]
Provisioning deis-controller with knife rackspace...
+ knife rackspace server create -y --server-create-timeout 1200 --server-name deis-controller --image 4b7c635d-89e1-44be-a15f-2877b5a660d1 --rackspace-region dfw --flavor 4 --identity-file /home/myuser/.ssh/deis-controller --bootstrap-version 11.4.4 --node-name deis-controller --run-list 'recipe[deis::controller]'
Instance ID: de17ca36-f186-4cdd-8969-4be58e7108ea
Name: deis-controller
Flavor: 2GB Standard Instance
...

$ ./contrib/rackspace/provision-rackspace-controller.sh dfw
...

```


[kniferack]: http://docs.opscode.com/plugin_knife_rackspace.html
