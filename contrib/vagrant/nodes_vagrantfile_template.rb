Vagrant.configure("2") do |config|
  config.vm.box = "deis-node"

  # The url from where the 'config.vm.box' box will be fetched if it
  # doesn't already exist on the user's system.
  config.vm.box_url = "https://s3-us-west-2.amazonaws.com/opdemand/deis-node.box"

  # Avahi-daemon will broadcast the node's address as $id.local
  config.vm.host_name = "$id"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network. IP will be fetched via DCHP and associated to '$id.local'
  # using avahi-daemon
  config.vm.network :public_network

  config.vm.provider :virtualbox do |vb|
    vb.customize ["modifyvm", :id, "--memory", "$memory"]
  end

  # Enable [hostname].local autodiscovery between VMs
  config.vm.provision :shell, inline: <<-SCRIPT
    # Avahi-daemon broadcasts the machine's hostname to local DNS.
    # So $id.local in this case
    sudo service avahi-daemon restart
  SCRIPT
end