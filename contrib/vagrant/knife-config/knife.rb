log_level                :info
log_location             STDOUT
node_name                'admin'
client_key               File.dirname(__FILE__) + '/../contrib/vagrant/knife-config/admin.pem'
validation_client_name   'chef-validator'
validation_key           File.dirname(__FILE__) + '/../contrib/vagrant/knife-config/chef-validator.pem'
chef_server_url          'https://chefserver.local'
syntax_check_cache_path  'syntax_check_cache'
