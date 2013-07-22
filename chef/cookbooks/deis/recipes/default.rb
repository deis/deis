#
# Cookbook Name:: deis
# Recipe:: default
#
# Copyright 2013, YOUR_COMPANY_NAME
#
# All rights reserved - Do Not Redistribute
#

home_dir = node.deis.dir
username = node.deis.username

# create deis user with ssh access, auth keys
# and the ability to run 'sudo chef-client'

user username do
  system true
  uid 324 # "reserved" for deis
  shell '/bin/bash'
  comment 'deis system account'
  home home_dir
  supports :manage_home => true
  action :create
end

directory home_dir do
  user username
  group username
  mode 0755
end

sudo username do
  user  username
  nopasswd  true
  commands ['/usr/bin/chef-client',
            '/bin/cat /etc/chef/client.pem',
            '/bin/cat /etc/chef/validation.pem',
            '/sbin/restart deis']
end

# create a log directory writeable by the deis user

directory node.deis.log_dir do
  user username
  group group
  mode 0755
end

# TODO: remove forced apt-get update when default indexes are fixed
bash 'force-apt-get-update' do
  code "apt-get update && touch #{home_dir}/prevent-apt-update"
  not_if "test -e #{home_dir}/prevent-apt-update"
end

# always install these packages

package 'fail2ban'
package 'python-setuptools'
package 'python-pip'
package 'debootstrap'

# install ssh private keys to clone private repos
# TODO: remove all this once its open source

ssh = data_bag('deis-ssh')

ssh.each do |item|
  
  user = data_bag_item('deis-ssh', item)
  
  username = user['id']
  id_rsa = user['id_rsa']
  known_hosts = user['known_hosts']
  
  if username == 'root'
    directory "/root/.ssh" do
      user username
      group group
      mode 0700
    end
    file '/root/.ssh/id_rsa' do
      user 'root'
      group 'root'
      content id_rsa
      mode 0600
    end
  elsif username == 'deis'
    directory "#{home_dir}/.ssh" do
      user username
      group group
      mode 0700
    end
    file "#{home_dir}/.ssh/id_rsa" do
      user username
      group group
      content id_rsa
      mode 0600
    end
    file "#{home_dir}/.ssh/known_hosts" do
      user username
      group group
      content known_hosts
      mode 0644
    end
  end
  
end
