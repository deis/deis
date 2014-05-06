# -*- mode: ruby -*-
# # vi: set ft=ruby :

# NOTE: This monkey-patching of the coreos guest plugin is a terrible
# hack that needs to be removed once the upstream plugin works with
# alpha CoreOS images.

require 'tempfile'
require 'ipaddr'
require Vagrant.source_root.join("plugins/guests/coreos/cap/configure_networks.rb")

BASE_CLOUD_CONFIG = <<EOF
#cloud-config

coreos:
    units:
      - name: coreos-cloudinit-vagrant-user.path
        command: start
        runtime: no
        content: |
          [Path]
          PathExists=/var/lib/coreos-vagrant/vagrantfile-user-data
      - name: coreos-cloudinit-vagrant-user.service
        runtime: no
        content: |
          [Unit]
          ConditionFileNotEmpty=/var/lib/coreos-vagrant/vagrantfile-user-data

          [Service]
          Type=oneshot
          EnvironmentFile=/etc/environment
          ExecStart=/usr/bin/coreos-cloudinit --from-file /var/lib/coreos-vagrant/vagrantfile-user-data
          RemainAfterExit=yes
EOF

NETWORK_UNIT = <<EOF
      - name: %s
        runtime: no
        content: |
          [Match]
          Name=%s

          [Network]
          Address=%s
EOF

# Borrowed from http://stackoverflow.com/questions/1825928/netmask-to-cidr-in-ruby
IPAddr.class_eval do
  def to_cidr
    self.to_i.to_s(2).count("1")
  end
end

module VagrantPlugins
  module GuestCoreOS
    module Cap
      class ConfigureNetworks
        include Vagrant::Util

        def self.configure_networks(machine, networks)
          cfg = BASE_CLOUD_CONFIG
          machine.communicate.tap do |comm|

            # Read network interface names
            interfaces = []
            comm.sudo("ifconfig | grep enp0 | cut -f1 -d:") do |_, result|
              interfaces = result.split("\n")
            end

            ip = ""

            # Configure interfaces
            # FIXME: fix matching of interfaces with IP adresses
            networks.each do |network|
              iface_num = network[:interface].to_i
              iface_name = interfaces[iface_num]
              cidr = IPAddr.new('255.255.255.0').to_cidr
              address = "%s/%s" % [network[:ip], cidr]
              unit_name = "50-%s.network" % [iface_name]
              unit = NETWORK_UNIT % [unit_name, iface_name, address]

              cfg = "#{cfg}#{unit}"
              ip = network[:ip]
            end

            cfg = <<EOF
#{cfg}
write_files:
  - path: /etc/environment
    content: |
      COREOS_PUBLIC_IPV4=#{ip}
      COREOS_PRIVATE_IPV4=#{ip}

hostname: #{machine.name}
EOF

            temp = Tempfile.new("coreos-vagrant")
            temp.write(cfg)
            temp.close

            comm.upload(temp.path, "/tmp/user-data")
            comm.sudo("mkdir -p /var/lib/coreos-vagrant")
            comm.sudo("mv /tmp/user-data /var/lib/coreos-vagrant/")

          end
        end
      end

      class ChangeHostName
        def self.change_host_name(machine, name)
          # This is handled in configure_networks
        end
      end
    end
  end
end
