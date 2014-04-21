# -*- mode: ruby -*-
# # vi: set ft=ruby :

require_relative 'contrib/coreos/override-plugin.rb'

Vagrant.configure("2") do |config|
  config.vm.box = "coreos-291.0.0"
  config.vm.box_url = "http://storage.core-os.net/coreos/amd64-usr/291.0.0/coreos_production_vagrant.box"

  config.vm.provider :vmware_fusion do |vb, override|
    override.vm.box_url = "http://storage.core-os.net/coreos/amd64-usr/291.0.0/coreos_production_vagrant_vmware_fusion.box"
  end

  config.vm.provider :virtualbox do |vb, override|
    vb.customize ["modifyvm", :id, "--memory", "4096"]
    # Fix docker not being able to resolve private registry in VirtualBox
    vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
    vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
  end

  # plugin conflict
  if Vagrant.has_plugin?("vagrant-vbguest") then
    config.vbguest.auto_update = false
  end

  config.vm.define vm_name = "deis" do |config|
    config.vm.hostname = vm_name

    ip = "172.17.8.100"
    config.vm.network :private_network, ip: ip

    # Uncomment below to enable NFS for sharing the host machine into the coreos-vagrant VM.
    config.vm.synced_folder ".", "/home/core/share", id: "core", :nfs => true, :mount_options => ['nolock,vers=3,udp']

    # disable update-engine to prevent reboots
    config.vm.provision :shell, :inline => "systemctl stop update-engine-reboot-manager && systemctl disable update-engine-reboot-manager && systemctl mask update-engine-reboot-manager", :privileged => true
    config.vm.provision :shell, :inline => "systemctl stop update-engine && systemctl disable update-engine && systemctl mask update-engine", :privileged => true

    # user-data bootstrapping
    config.vm.provision :file, :source => "contrib/coreos/user-data", :destination => "/tmp/vagrantfile-user-data"
    config.vm.provision :shell, :inline => "mv /tmp/vagrantfile-user-data /var/lib/coreos-vagrant/", :privileged => true

  end

end
