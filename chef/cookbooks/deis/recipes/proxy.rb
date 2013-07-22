#
# Cookbook Name:: deis-proxy
# Recipe:: default
#
# Copyright 2013, YOUR_COMPANY_NAME
#
# All rights reserved - Do Not Redistribute
#

include_recipe 'deis::nginx'

# iterate over this node's formation

formations = data_bag('deis-formations')

formations.each do |f|
  
  formation = data_bag_item('deis-formations', f)
  
  # skip this node if it's not configured as a proxy
  next if ! formation['nodes']['proxies'].keys.include? node.name
  
  proxy = formation['proxy']
  
  vars = {:formation => f,
          :port => proxy['port'],
          :backends => proxy['backends'], 
          :algorithm => proxy['algorithm'],
          :firewall => proxy['firewall']}

  nginx_site "deis-#{f}" do
    template 'nginx-proxy.conf.erb'
    vars vars
  end
  
end
