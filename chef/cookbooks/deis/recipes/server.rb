
username = node.deis.username
group = node.deis.group
controller_dir = node.deis.controller.dir

# required packages

package 'python-virtualenv'
package 'python-dev'
package 'libevent-dev' # for gevent
package 'rabbitmq-server' # for celery
package 'libpq-dev' # for psycopg2

# synchronize the gitosis repository

git controller_dir do
  user username
  group group
  repository node.deis.controller.repository
  action :sync
end

directory controller_dir do
  user username
  group group
  mode 0755 # need nginx access to static files
end

# write out local settings for db access, etc.

template "#{controller_dir}/deis/local_settings.py" do
  user username
  group group
  mode 0644
  source 'local_settings.py.erb'
  variables :debug => node.deis.controller.debug, 
            :secret_key => node.deis.controller.secret_key,
            :db_name => node.deis.database.name,
            :db_user => node.deis.database.user
  subscribes :create, "git[#{controller_dir}]", :immediately
  notifies :restart, "service[deis-server]", :delayed
  notifies :restart, "service[deis-worker]", :delayed
end

# virtualenv setup

bash 'deis-controller-virtualenv' do
  user username
  group group
  cwd controller_dir
  code "virtualenv --distribute venv"
  creates "#{controller_dir}/venv"
  action :nothing
  subscribes :run, "git[#{controller_dir}]", :immediately 
end

bash 'deis-controller-pip-install' do
  user username
  group group
  cwd controller_dir
  code "source venv/bin/activate && pip install -r requirements.txt"
  action :nothing
  subscribes :run, "git[#{controller_dir}]", :immediately 
end

# NOTE: collectstatic and other subcommands must be run after local_settings

bash 'deis-controller-collectstatic' do
  user username
  group group
  cwd controller_dir
  code "source venv/bin/activate && ./manage.py collectstatic --noinput"
  action :nothing
  subscribes :run, "git[#{controller_dir}]", :immediately 
end

# write out upstart daemon

template '/etc/init/deis-server.conf' do
  user 'root'
  group 'root'
  mode 0644
  source 'deis-server.conf.erb'
  variables :home => node.deis.controller.dir,
            :port => node.deis.controller.worker_port,
            :bind => '0.0.0.0',
            :workers => node.deis.controller.workers
  notifies :restart, "service[deis-server]", :delayed
end

service 'deis-server' do
  provider Chef::Provider::Service::Upstart
  action [:enable]
  subscribes :restart, "git[#{controller_dir}]", :delayed
end

template '/etc/init/deis-worker.conf' do
  user 'root'
  group 'root'
  mode 0644
  source 'deis-worker.conf.erb'
  variables :home => node.deis.controller.dir
  notifies :restart, "service[deis-worker]", :delayed            
end

service 'deis-worker' do
  provider Chef::Provider::Service::Upstart
  action [:enable]
  subscribes :restart, "git[#{controller_dir}]", :delayed
end

# nginx configuration

include_recipe 'deis::nginx'

nginx_site 'deis-controller' do
  template 'nginx-controller.conf.erb'
  vars :server_root => node.deis.controller.dir,
       :slug_root => node.deis.build.slug_dir,
       :http_port => node.deis.controller.http_port
end

