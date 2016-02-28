# -*- mode: ruby -*-
# # vi: set ft=ruby :

require 'fileutils'
require 'open3'

Vagrant.require_version ">= 1.6.5"

unless Vagrant.has_plugin?("vagrant-triggers")
  raise Vagrant::Errors::VagrantError.new, "Please install the vagrant-triggers plugin running 'vagrant plugin install vagrant-triggers'"
end

CLOUD_CONFIG_PATH = File.join(File.dirname(__FILE__), "contrib", "coreos", "user-data")
CONFIG = File.join(File.dirname(__FILE__), "config.rb")
CONTRIB_UTILS_PATH = File.join(File.dirname(__FILE__), "contrib", "utils.sh")

# Make variables from contrib/utils.sh accessible
if File.exists?(CONTRIB_UTILS_PATH)
  cu_vars = Hash.new do |hash, key|
    stdin, stdout, stderr = Open3.popen3("/usr/bin/env", "bash", "-c", "source '#{CONTRIB_UTILS_PATH}' && echo $#{key}")
    value = stdout.gets.chomp
    hash[key] = value unless value.empty?
  end
else
  raise Vagrant::Errors::VagrantError.new, "The file '#{CONTRIB_UTILS_PATH}' is missing."
end

# Defaults for config options defined in CONFIG
$num_instances = 1
$instance_name_prefix = "deis"
$update_channel = cu_vars["COREOS_CHANNEL"]
$image_version = cu_vars["COREOS_VERSION"]
$enable_serial_logging = false
$share_home = false
$vm_gui = false
$vm_memory = 2048
$vm_cpus = 1
$shared_folders = {}
$forwarded_ports = {}

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

# Use old vb_xxx config variables when set
def vm_gui
  $vb_gui.nil? ? $vm_gui : $vb_gui
end

def vm_memory
  $vb_memory.nil? ? $vm_memory : $vb_memory
end

def vm_cpus
  $vb_cpus.nil? ? $vm_cpus : $vb_cpus
end

Vagrant.configure("2") do |config|
  # always use Vagrants insecure key
  config.ssh.insert_key = false

  config.vm.box = "coreos-%s" % $update_channel
  if $image_version != "current"
      config.vm.box_version = $image_version
  end
  config.vm.box_url = "http://%s.release.core-os.net/amd64-usr/%s/coreos_production_vagrant.json" % [$update_channel, $image_version]

  ["vmware_fusion", "vmware_workstation"].each do |vmware|
    config.vm.provider vmware do |v, override|
      override.vm.box_url = "http://%s.release.core-os.net/amd64-usr/%s/coreos_production_vagrant_vmware_fusion.json" % [$update_channel, $image_version]
    end
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
    if File.exists?(CLOUD_CONFIG_PATH) && !File.readlines(CLOUD_CONFIG_PATH).grep(/\s*discovery #DISCOVERY_URL/).any?
      user_data = File.read(CLOUD_CONFIG_PATH)
      new_userdata = user_data.gsub("coreos:", "coreos:\n  flannel:\n    interface: $public_ipv4")
      File.open(CLOUD_CONFIG_PATH, "w") {|file| file.puts new_userdata }
    else
      raise Vagrant::Errors::VagrantError.new, "Run 'make discovery-url' first to create user-data."
    end
  end

  (1..$num_instances).each do |i|
    config.vm.define vm_name = "%s-%02d" % [$instance_name_prefix, i] do |config|
      config.vm.hostname = vm_name

      if $enable_serial_logging
        logdir = File.join(File.dirname(__FILE__), "log")
        FileUtils.mkdir_p(logdir)

        serialFile = File.join(logdir, "%s-serial.txt" % vm_name)
        FileUtils.touch(serialFile)

        ["vmware_fusion", "vmware_workstation"].each do |vmware|
          config.vm.provider vmware do |v, override|
            v.vmx["serial0.present"] = "TRUE"
            v.vmx["serial0.fileType"] = "file"
            v.vmx["serial0.fileName"] = serialFile
            v.vmx["serial0.tryNoRxLoss"] = "FALSE"
          end
        end

        config.vm.provider :virtualbox do |vb, override|
          vb.customize ["modifyvm", :id, "--uart1", "0x3F8", "4"]
          vb.customize ["modifyvm", :id, "--uartmode1", serialFile]
        end
      end

      if $expose_docker_tcp
        config.vm.network "forwarded_port", guest: 2375, host: ($expose_docker_tcp + i - 1), auto_correct: true
      end

      $forwarded_ports.each do |guest, host|
        config.vm.network "forwarded_port", guest: guest, host: host, auto_correct: true
      end

      ["vmware_fusion", "vmware_workstation"].each do |vmware|
        config.vm.provider vmware do |v|
          v.gui = vm_gui
          v.vmx['memsize'] = vm_memory
          v.vmx['numvcpus'] = vm_cpus
        end
      end

      config.vm.provider :virtualbox do |vb|
        vb.gui = vm_gui
        vb.memory = vm_memory
        vb.cpus = vm_cpus
      end

      ip = "172.17.8.#{i+99}"
      config.vm.network :private_network, ip: ip

      # Use the same nameserver as the host machine in order to avoid the "too many redirects" problem.
      config.vm.provider :virtualbox do |vb|
        vb.customize ["modifyvm", :id, "--natdnshostresolver1", "off"]
        vb.customize ["modifyvm", :id, "--natdnsproxy1", "off"]
      end

      # Uncomment below to enable NFS for sharing the host machine into the coreos-vagrant VM.
      #config.vm.synced_folder ".", "/home/core/share", id: "core", :nfs => true, :mount_options => ['nolock,vers=3,udp']
      $shared_folders.each_with_index do |(host_folder, guest_folder), index|
        config.vm.synced_folder host_folder.to_s, guest_folder.to_s, id: "core-share%02d" % index, nfs: true, mount_options: ['nolock,vers=3,udp']
      end

      if $share_home
        config.vm.synced_folder ENV['HOME'], ENV['HOME'], id: "home", :nfs => true, :mount_options => ['nolock,vers=3,udp']
      end

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
