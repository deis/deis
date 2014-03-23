Vagrant.configure("2") do |config|
  config.vm.box = "deis-server"

  # The url from where the 'config.vm.box' box will be fetched if it
  # doesn't already exist on the user's system.
  config.vm.box_url = "https://s3-us-west-2.amazonaws.com/opdemand/ubuntu-12.04.3-amd64-vbox.box"

  # Avahi-daemon will broadcast the node's address as $id.local
  config.vm.host_name = "$id"

  # IP will be associated to '$id.local' using avahi-daemon
  config.vm.network :private_network, ip: "$ipaddress"

  config.vm.provider :virtualbox do |vb|
    vb.customize ["modifyvm", :id, "--memory", "$memory"]
  end

  # Enable [hostname].local autodiscovery between VMs
  config.vm.provision :shell, inline: <<-SCRIPT
    # Avahi-daemon broadcasts the machine's hostname to local DNS.
    # So $id.local in this case
    sudo apt-get install -yq avahi-daemon
    sudo service avahi-daemon restart
    echo "*.* @@192.168.61.100:514" > /etc/rsyslog.d/51-remote.conf
  SCRIPT
end
