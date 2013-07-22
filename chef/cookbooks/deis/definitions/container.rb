
define :container, :c_type => nil, :c_num => nil, :env => {}, :command => '', :port => nil, :image => nil, :slug_dir => nil, :enable => nil, :user => "root" do
  
  # pull out local variables
  c_type = params[:c_type]
  c_num = params[:c_num]
  env = params[:env]
  command = params[:command]
  port = params[:port]
  image = params[:image]
  slug_dir = params[:slug_dir]
  user = params[:user]
  
  # see if service should be enabled or disabled
  if params[:enable] == true
    actions = [:start, :enable]
    on_change = :restart
  else
    actions = [:stop, :disable]
    on_change = :nothing
  end
  
  # create upstart service definition
  template "/etc/init/#{c_type}.#{c_num}.conf" do
    source "container.conf.erb"
    mode 0644
    variables({
      :log_dir => "/var/log",
      :user => user,
      :image => image,
      :slug_dir => slug_dir,
      :env => env,
      :port => port,
      :command => command,     
      :c_type => c_type,
      :c_num => c_num
    })
    # stop the service to force job definition reload
    notifies :stop, "service[#{c_type}.#{c_num}]", :immediately
    notifies on_change, "service[#{c_type}.#{c_num}]", :delayed
  end
  
  # define an upstart daemon as enabled or disabled
  service "#{c_type}.#{c_num}" do
    provider Chef::Provider::Service::Upstart
    action actions
  end

end