
include_recipe 'apt'

apt_repository 'docker-ppa' do
  uri 'http://ppa.launchpad.net/dotcloud/lxc-docker/ubuntu'
  distribution node['lsb']['codename']
  components ['main']
  keyserver 'keyserver.ubuntu.com'
  key '63561DC6'
end

package 'lxc-docker'

service 'docker' do
  provider Chef::Provider::Service::Upstart  
  supports :status => true, :restart => true, :reload => true
  action [ :enable ]
end
