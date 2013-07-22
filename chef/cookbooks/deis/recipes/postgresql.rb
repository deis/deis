package 'postgresql-9.1'

service 'postgresql' do
  supports :status => true, :restart => true, :reload => true
  action [ :enable, :start ]
end

template '/etc/postgresql/9.1/main/pg_hba.conf' do
  source 'pg_hba.conf.erb'
  user 'postgres'
  group 'postgres'
  mode 0640
  notifies :reload, resources(:service => 'postgresql')
end 

template '/etc/postgresql/9.1/main/postgresql.conf' do
  source 'postgresql.conf.erb'
  user 'postgres'
  group 'postgres'
  mode 0644
  notifies :reload, resources(:service => 'postgresql')
end

db_name = node.deis.database.name
db_user = node.deis.database.user

execute 'create-deis-database' do
    user 'postgres'
    group 'postgres'
    db_exists = <<-EOF
    psql -c "select * from pg_database WHERE datname='#{db_name}'" | grep -c #{db_name}
    EOF
    command "createdb --encoding=utf8 --template=template0 #{db_name}"
    not_if db_exists, :user => 'postgres'
end

execute 'create-deis-database-user' do
    user 'postgres'
    group 'postgres'
    user_exists = <<-EOF
    psql -c "select * from pg_user where usename='#{db_user}'" | grep -c #{db_user}
    EOF
    command "createuser --no-superuser --no-createrole --no-createdb --no-password #{db_user}"
    not_if user_exists, :user => 'postgres'
end
