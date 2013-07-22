Creating SSL certificates is a common task done in web application infrastructures, so a rake task is provided to generate certificates.  These certificates are stored here by the ssl_cert task.  

Configure the values used in the SSL certificate by modifying `config/rake.rb`.

To generate a certificate set for a new monitoring server, for example:

    rake ssl_cert FQDN=monitoring.example.com

Once the certificates are generated, copy them into the cookbook(s) where you want to use them.

    cp certificates/monitoring.example.com.* cookbooks/COOKBOOK/files/default

In the recipe for that cookbook, create a `cookbook_file` resource to configure a resource that puts them in place on the destination server.

    cookbook_file '/etc/apache2/ssl/monitoring.example.com.pem'
      owner 'root'
      group 'root'
      mode 0600
    end
