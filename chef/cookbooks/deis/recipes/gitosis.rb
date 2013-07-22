#
# Cookbook Name:: gitosis
# Recipe:: default
#
# Copyright 2013, deis
#
# All rights reserved - Do Not Redistribute
#

gitosis_dir = node.deis.gitosis.dir
gitosis_checkout = "#{gitosis_dir}/checkout"
gitosis_admin_repo = "#{gitosis_dir}/repositories/gitosis-admin.git"
gitosis_admin_checkout = "#{gitosis_dir}/gitosis-admin"
gitosis_key_dir = "#{gitosis_admin_checkout}/keydir"

# setup git core

package 'git'

# create git user

user 'git' do
  system true
  uid 325 # "reserved" for git
  shell '/bin/sh'
  comment 'git version control'
  home gitosis_dir
  supports :manage_home => true
  action :create
end

# allow git user to trigger build-release-run script during
# git push build hook

sudo 'git' do
  user  'git'
  runas node.deis.username
  nopasswd  true
  commands [ node.deis.controller.dir + '/bin/build-release-run' ]
end

# synchronize the gitosis repository

git gitosis_checkout do
  repository 'git://github.com/opdemand/gitosis.git'
  action :sync
end

# install gitosis

bash "gitosis-install" do
  cwd gitosis_checkout
  code "python setup.py install"
  action :nothing
  subscribes :run, "git[#{gitosis_checkout}]", :immediately
end

# initialize gitosis

bash 'gitosis-generate-ssh-key' do
  user 'git'
  group 'git'
  code "ssh-keygen -t rsa -b 2048 -N '' -f #{gitosis_dir}/.ssh/gitosis-admin -C gitosis-admin"
  not_if "test -e #{gitosis_dir}/.ssh/gitosis-admin"
end

bash 'gitosis-init' do
  code "sudo -H -u git gitosis-init < #{gitosis_dir}/.ssh/gitosis-admin.pub"
  not_if "test -e #{gitosis_dir}/gitosis-admin"
end

bash 'git-clone-gitosis-admin' do
  user 'git'
  group 'git'
  cwd gitosis_dir
  code "git clone #{gitosis_admin_repo}"
  not_if "test -e #{gitosis_admin_checkout}"
end

# try to load the gitosis data bag item

gitosis = data_bag_item('deis-build', 'gitosis')

# create ssh keys

gitosis['ssh_keys'].each do |key_name, key_material|
  file "#{gitosis_key_dir}/#{key_name}.pub" do
    owner 'git'
    group 'git'
    content key_material
    notifies :run, 'bash[git-add-gitosis-admin]'
  end
end

# purge old ssh keys

Dir.glob("#{gitosis_key_dir}/*.pub").each do |f|
  next if f.sub("#{gitosis_key_dir}/", '') == 'gitosis-admin.pub'
  if !gitosis['ssh_keys'].has_key? f.sub(gitosis_key_dir+'/', '').sub('.pub', '')
    file f do
      action :delete
      notifies :run, 'bash[git-add-gitosis-admin]'
    end
  end
end

# configure gitosis

template "#{gitosis_admin_checkout}/gitosis.conf" do
  user 'git'
  group 'git'
  source 'gitosis.conf.erb'
  variables({
    :admins => ['gitosis-admin'],
    :formations => gitosis['formations'],
  })
  notifies :run, 'bash[git-add-gitosis-admin]'
end

bash 'git-add-gitosis-admin' do
  user 'git'
  group 'git'
  cwd gitosis_admin_checkout
  code 'git add .'
  action :nothing
  notifies :run, 'bash[git-commit-gitosis-admin]'
end

bash 'git-commit-gitosis-admin' do
  user 'git'
  group 'git'
  cwd gitosis_admin_checkout
  code 'git commit -m "auto-update via chef"'
  action :nothing
  notifies :run, 'bash[git-push-gitosis-admin]'
end

bash 'git-push-gitosis-admin' do
  user 'git'
  group 'git'
  cwd gitosis_admin_checkout
  code 'git push'
  action :nothing
  notifies :run, 'bash[gitosis-update-hook]'
end

# TODO: figure out why this needs to be run manually
# after pushing to the gitosis-admin repo
bash 'gitosis-update-hook' do
  user 'git'
  group 'git'
  cwd gitosis_admin_repo
  code 'gitosis-run-hook post-update'
  environment 'GIT_DIR' => gitosis_admin_repo, 'HOME' => gitosis_dir
  action :nothing
end

# create application repositories

gitosis['formations'].each_pair do |name, ssh_keys|
  dir = "#{gitosis_dir}/repositories/#{name}.git"
  directory dir do
    user 'git'
    group 'git'
    mode 0750
  end
  bash 'gitosis-init-bare-app' do
    user 'git'
    group 'git'
    cwd dir
    code 'git init --bare'
    not_if "test -e #{dir}/HEAD"
  end
end

# purge old repository directories

Dir.glob("#{gitosis_dir}/repositories/*.git").each do |d|
  next if d.include? 'gitosis-admin.git'
  if !gitosis['formations'].has_key? d.sub("#{gitosis_dir}/repositories/", '').sub('.git', '')
    directory d do
      action :delete
      recursive true
    end
  end
end
