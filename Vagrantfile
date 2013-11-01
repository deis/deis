
Vagrant.configure("2") do |config|
  config.vm.box = "deis-node"

  # This is the vanilla Ubunutu 12.04 Precise box. It's about 350MB
  config.vm.box_url = "https://s3-us-west-2.amazonaws.com/opdemand/deis-node.box"

  # This is a premade Deis controller box, probably 12.04, need to check. It's about 1.1GB
  # config.vm.box_url = "https://s3-us-west-2.amazonaws.com/opdemand/deis-controller.box"

  # Avahi-daemon will broadcast the server's address as deis-controller.local
  config.vm.host_name = "deis-controller"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network. IP will be fetched via DCHP and associated to 'chefserver.local'
  # using avahi-daemon
  config.vm.network :public_network

  # Chef Server requires at least 1G of RAM to install.
  # You may be able to run it with less once it's installed.
  config.vm.provider :virtualbox do |vb|
    vb.customize ["modifyvm", :id, "--memory", "1024"]
  end

  # 'deis provider:discover' detects the host machine's user and IP address, however, that command cannot
  # be guareteed to run inside the deis codebase. Therefore we can't use that opportunity to discover
  # the path of the codebase on the host machine. Therefore we do it now as this Vagrantfile has to exist
  # inside the codebase.
  nodes_dir = File.dirname(__FILE__) + '/contrib/vagrant/nodes'

  config.vm.provision :shell, inline: <<-SCRIPT
    # Avahi-daemon broadcasts the machine's hostname to local DNS.
    # Therefore 'deis-controller.local' in this case.
    sudo service avahi-daemon restart
    # Make a record of where the deis code base is on the host machine
    echo "#{nodes_dir}" > /home/vagrant/.host_nodes_dir
  SCRIPT
end

# If you want to do some funky custom stuff to your box, but don't want those things tracked by git,
# add a Vagrantfile.local and it will be included. You can use the exact same syntax as above. For
# example you could mount your dev version of deis onto the VM and hack live on the VM;
# `config.vm.share_folder "deis", "/opt/deis", "~/myworkspace/deis"
# Or if you're low on RAM you can boot the VM with less RAM. Note that at least 1GB is needed for
# installation, but you may be able to get away with 512MB once everything is installed.
load "Vagrantfile.local" if File.exists? "Vagrantfile.local"
