# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

  config.vm.box = "deis-controller"
  config.vm.hostname = "deis-controller"

  config.vm.provider :virtualbox do |v|
    v.customize ["modifyvm", :id, "--memory", 2048]
  end
  
  config.vm.box_url = "http://files.vagrantup.com/precise64.box"

  config.vm.network :public_network, :bridge => 'en0: Wi-Fi (AirPort)' #, :mac => "08002769c9a0"  

  config.vm.provision :shell, :inline => "echo Bootstrap with: knife bootstrap `/sbin/ifconfig eth1|grep inet|head -1|sed 's/\:/ /'|awk '{print $3}'` -x vagrant -P vagrant -N deis-controller -r role[deis-controller] --sudo"
  
end
