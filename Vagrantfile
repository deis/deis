# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  # Give each controller and node additional memory
  config.vm.provider :virtualbox do |v|
    v.customize ["modifyvm", :id, "--memory", 2048]
  end

  # Deis Nodes
  config.vm.box_url = "https://s3-us-west-2.amazonaws.com/opdemand/deis-node.box"
  config.vm.box = "deis-node"

  # Deis Controller
  config.vm.define "deis-controller", primary: true do |controller|
    controller.vm.hostname = "deis-controller"
    controller.vm.network "private_network", ip: "192.168.61.100"
  end

  # Node 1
  config.vm.define "deis-node-1" do |node1|
    node1.vm.hostname = "deis-node-1"
    node1.vm.network "private_network", ip: "192.168.61.101"
  end

  # Node 2
  config.vm.define "deis-node-2" do |node2|
    node2.vm.hostname = "deis-node-2"
    node2.vm.network "private_network", ip: "192.168.61.102"
  end

end
