# -*- mode: ruby -*-
# # vi: set ft=ruby :

# <hack>

# NOTE: This monkey-patching of the coreos guest plugin is a terrible
# hack that needs to be removed once the upstream plugin works with the
# new CoreOS images.

require Vagrant.source_root.join("plugins/guests/coreos/cap/configure_networks.rb")

module VagrantPlugins
  module GuestCoreOS
	module Cap
	  class ConfigureNetworks
		include Vagrant::Util

		def self.configure_networks(machine, networks)
		  machine.communicate.tap do |comm|
			# Read network interface names
			interfaces = []
			comm.sudo("ifconfig | grep enp0 | cut -f1 -d:") do |_, result|
			  interfaces = result.split("\n")
			end

			# Configure interfaces
			# FIXME: fix matching of interfaces with IP adresses
			networks.each do |network|
			  comm.sudo("ifconfig #{interfaces[network[:interface].to_i]} #{network[:ip]} netmask #{network[:netmask]}")
			end

		  end
		end
	  end

      class ChangeHostName
        def self.change_host_name(machine, name)
			# do nothing!
        end
      end
    end
  end
end

# </hack>
