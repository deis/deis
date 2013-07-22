
username = node.deis.username
group = node.deis.group
client_dir = node.deis.client.dir

git client_dir do
  user username
  group group
  repository node.deis.client.repository
  action :sync
end

directory client_dir do
  user username
  group group
  mode 0755
end

bash 'deis-client-pip-install' do
  cwd client_dir
  code "pip install -r requirements.txt"
  action :nothing
  subscribes :run, "git[#{client_dir}]", :immediately
end

bash 'deis-client-python-install' do
  cwd client_dir
  code "python setup.py install"
  action :nothing
  subscribes :run, "git[#{client_dir}]", :immediately  
end