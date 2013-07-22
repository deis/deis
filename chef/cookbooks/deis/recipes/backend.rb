include_recipe 'deis::docker'

username = node.deis.username
group = node.deis.group
home = node.deis.dir
image = node.deis.runtime.image

directory node.deis.runtime.dir do
  user username
  group group
  mode 0700
end

directory node.deis.runtime.slug_root do
  user username
  group group
  mode 0700
end

package 'curl'

formations = data_bag('deis-formations')

formations.each do |f|
  
  formation = data_bag_item('deis-formations', f)

  # skip this node if it's not configured as a proxy
  next if ! formation['nodes']['backends'].keys.include? node.name
  
  id = formation['id']
  version = formation['release']['version']
  build = formation['release']['build']
  config = formation['release']['config']
  image = formation['release']['image']
  
  # pull the image if it doesn't exist already
  
  bash "pull-image-#{image}" do
    cwd node.deis.runtime.dir
    code "docker pull #{image}"
    not_if "docker images | grep #{image}"
  end
  
  # if build is specified, use special heroku-style runtime
  
  
  if build != {}
  
    slug_url = build['url']
    
    # download the slug to a tempdir
    slug_root = node.deis.runtime.slug_root
    slug_dir = "#{slug_root}/#{f}-#{version}"
    slug_filename = "app.tar.gz"
    slug_path = "#{slug_dir}/#{slug_filename}"
    
    bash "download-slug-#{f}-#{version}" do
      cwd slug_root
      code <<-EOF
        rm -rf #{slug_dir}
        mkdir -p #{slug_dir}
        cd #{slug_dir}
        curl -s #{slug_url} > #{slug_path}
        tar xvfz #{slug_path}
        EOF
      not_if "test -f #{slug_path}"
    end
  else
    slug_dir = nil
  end

  # iterate over this application's process formation by
  # Procfile-defined type
  
  formation['containers'].each_pair do |c_type, c_formation|
    
    c_formation.each_pair do |c_num, node_port|
    
      nodename, port = node_port.split(':')
      
      # if the nodename doesn't match don't enable the process
      # but still define it and leave it disabled
      if nodename == node.name
        enabled = true
      else
        enabled = false
      end
      
      # determine build command, if one exists
      if build != {}
        command = build['procfile'][c_type]
      else
        command = nil # assume command baked into docker image
      end
      
      # define the container
      container "#{c_type}.#{c_num}" do
        c_type c_type
        c_num c_num
        env config
        command command
        port port
        image image
        slug_dir slug_dir
        enable enabled
        user username
      end      
  
    end
    
    # cleanup any old containers that match this process type
    (1..100).each { |n|
      
      # skip this c_num if we already processed it
      unless c_formation.has_key?(n.to_s)
        filename = "/etc/init/#{c_type}.#{n}.conf"
        
        # see if the upstart service exists
        if File.exist?(filename) or File.exist?(filename+".old")
        
          # stop and disable it
          service "#{c_type}.#{n}" do
            provider Chef::Provider::Service::Upstart
            action [:stop, :disable]
          end
          
          # delete the service definition and any *.old files
          [ filename, "#{filename}.old"].each { |fl|
              file fl do
                action :delete
              end
            }
          end
        end
      }
    
  end

end

