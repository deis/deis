name "deis-controller"
description "Deis PaaS Controller"
run_list "recipe[deis]", "recipe[deis::gitosis]", "recipe[deis::build]", "recipe[deis::postgresql]", "recipe[deis::server]", "recipe[deis::client]"
#env_run_lists "prod" => ["recipe[apache2]"], "staging" => ["recipe[apache2::staging]"], "_default" => []
#default_attributes "apache2" => { "listen_ports" => [ "80", "443" ] }
#override_attributes "apache2" => { "max_children" => "50" }
