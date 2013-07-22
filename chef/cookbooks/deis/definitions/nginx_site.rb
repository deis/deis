
define :nginx_site, :name => nil, :template => nil, :vars => {} do

  name = params[:name]
  template = params[:template]
  vars = params[:vars]
  
  template "/etc/nginx/sites-available/#{name}" do
    source template
    mode 0644
    variables(vars)
    # TODO: switch to no-downtime reload once listen port changes
    # are respected, for now port changes are ignored
    notifies :restart, "service[nginx]", :delayed
  end
  
  link "/etc/nginx/sites-enabled/#{name}" do
    to "/etc/nginx/sites-available/#{name}"
    action :create
  end

end
