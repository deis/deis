name "deis-backend"
description "Deis PaaS Backend"
run_list "recipe[deis]", "recipe[deis::backend]"
#env_run_lists "prod" => ["recipe[apache2]"], "staging" => ["recipe[apache2::staging]"], "_default" => []
#default_attributes "apache2" => { "listen_ports" => [ "80", "443" ] }
#override_attributes "apache2" => { "max_children" => "50" }
