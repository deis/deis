# -*- mode: ruby -*-
# # vi: set ft=ruby :

require 'fileutils'

Vagrant.require_version ">= 1.6.5"

unless Vagrant.has_plugin?("vagrant-triggers")
  raise Vagrant::Errors::VagrantError.new, "Please install the vagrant-triggers plugin running 'vagrant plugin install vagrant-triggers'"
end

CLOUD_CONFIG_PATH = File.join(File.dirname(__FILE__), "contrib", "coreos", "user-data")
CONFIG = File.join(File.dirname(__FILE__), "config.rb")

# Defaults for config options defined in CONFIG
$num_instances = 1
$update_channel = ENV["COREOS_CHANNEL"] || "stable"
$enable_serial_logging = false
$vb_gui = false
$vb_memory = 2048
$vb_cpus = 1

# Attempt to apply the deprecated environment variable NUM_INSTANCES to
# $num_instances while allowing config.rb to override it
if ENV["NUM_INSTANCES"].to_i > 0 && ENV["NUM_INSTANCES"]
  $num_instances = ENV["NUM_INSTANCES"].to_i
elsif ENV["DEIS_NUM_INSTANCES"].to_i > 0 && ENV["DEIS_NUM_INSTANCES"]
  $num_instances = ENV["DEIS_NUM_INSTANCES"].to_i
else
  $num_instances = 3
end

if File.exist?(CONFIG)
  require CONFIG
end

Vagrant.configure("2") do |config|
  # always use Vagrants insecure key
  config.ssh.insert_key = false

  config.vm.box = "coreos-%s" % $update_channel
  config.vm.box_version = ">= 522.6.0"
  config.vm.box_url = "http://%s.release.core-os.net/amd64-usr/current/coreos_production_vagrant.json" % $update_channel

  config.vm.provider :vmware_fusion do |vb, override|
    override.vm.box_url = "http://%s.release.core-os.net/amd64-usr/current/coreos_production_vagrant_vmware_fusion.json" % $update_channel
  end

  config.vm.provider :virtualbox do |v|
    # On VirtualBox, we don't have guest additions or a functional vboxsf
    # in CoreOS, so tell Vagrant that so it can be smarter.
    v.check_guest_additions = false
    v.functional_vboxsf     = false
  end

  # plugin conflict
  if Vagrant.has_plugin?("vagrant-vbguest") then
    config.vbguest.auto_update = false
  end

  config.trigger.before :up do
    if !File.exists?(CLOUD_CONFIG_PATH) || File.readlines(CLOUD_CONFIG_PATH).grep(/#\s*discovery:/).any?
      raise Vagrant::Errors::VagrantError.new, "Run 'make discovery-url' first to create user-data."
    end
  end

  (1..$num_instances).each do |i|
    config.vm.define vm_name = "deis-%d" % i do |config|
      config.vm.hostname = vm_name

      if $enable_serial_logging
        logdir = File.join(File.dirname(__FILE__), "log")
        FileUtils.mkdir_p(logdir)

        serialFile = File.join(logdir, "%s-serial.txt" % vm_name)
        FileUtils.touch(serialFile)

        config.vm.provider :vmware_fusion do |v, override|
          v.vmx["serial0.present"] = "TRUE"
          v.vmx["serial0.fileType"] = "file"
          v.vmx["serial0.fileName"] = serialFile
          v.vmx["serial0.tryNoRxLoss"] = "FALSE"
        end

        config.vm.provider :virtualbox do |vb, override|
          vb.customize ["modifyvm", :id, "--uart1", "0x3F8", "4"]
          vb.customize ["modifyvm", :id, "--uartmode1", serialFile]
        end
      end

      if $expose_docker_tcp
        config.vm.network "forwarded_port", guest: 2375, host: ($expose_docker_tcp + i - 1), auto_correct: true
      end

      config.vm.provider :vmware_fusion do |vb|
        vb.gui = $vb_gui
      end

      config.vm.provider :virtualbox do |vb|
        vb.gui = $vb_gui
        vb.memory = $vb_memory
        vb.cpus = $vb_cpus
      end

      ip = "172.17.8.#{i+99}"
      config.vm.network :private_network, ip: ip

      # Uncomment below to enable NFS for sharing the host machine into the coreos-vagrant VM.
      #config.vm.synced_folder ".", "/home/core/share", id: "core", :nfs => true, :mount_options => ['nolock,vers=3,udp']

      if File.exist?(CLOUD_CONFIG_PATH)
        config.vm.provision :file, :source => "#{CLOUD_CONFIG_PATH}", :destination => "/tmp/vagrantfile-user-data"
        # check that the CoreOS user-data file is valid
        config.vm.provision :shell do |s|
          s.path = File.join(File.dirname(__FILE__), "contrib", "util", "check-user-data.sh")
          s.args = ["/tmp/vagrantfile-user-data", $num_instances]
        end
        config.vm.provision :shell, :inline => "mv /tmp/vagrantfile-user-data /var/lib/coreos-vagrant/", :privileged => true
      else
        config.vm.provision :shell do |s|
          s.inline = "echo \"File not found: #{CLOUD_CONFIG_PATH}\" &&" +
            "echo \"Run 'make discovery-url' first to create user-data.\" && exit 1"
        end
      end

    end
  end
end
