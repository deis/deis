deis-chef
=========

**Deis** is an open source *platform-as-a-service* (PaaS) for public and
private clouds.

Take your agile development to the next level. Free your mind and focus on
your code. Deploy updates to local metal or worldwide clouds with
```git push```. Scale servers, processes, and proxies with a simple command.
Enjoy the *twelve-factor app* workflow while keeping total control.

The [opdemand/deis-chef](https://github.com/opdemand/deis-chef) project
contains OpsCode Chef recipes for provisioning servers you create in **Deis**.
To set up your own private application platform, the
[deis-controller](https://github.com/opdemand/deis-controller) and
[deis](https://github.com/opdemand/deis) projects are also required.


Getting Started
---------------

How to configure the OpsCode Chef components for **Deis** v0.0.2. The **Deis**
server/API component is in the deis-controller project.

##### Clone the repository

    git clone git@github.com:opdemand/deis-chef.git

##### Create Chef Organization

* Browse to <http://www.opscode.com>
* Sign up for a new account/organization
* Download the account key to `~/.chef/gabrtv.pem`
* Download the validator key to `~/.chef/gabrtv-validator.pem`

##### Set up Knife

Setup `~/.chef/knife.rb` using the following template

	current_dir = File.dirname(__FILE__)
	log_level                :info
	log_location             STDOUT
	node_name                "gabrtv"
	client_key               "#{current_dir}/gabrtv.pem"
	validation_client_name   "gabrtv-validator"
	validation_key           "#{current_dir}/gabrtv-validator.pem"
	chef_server_url          "https://api.opscode.com/organizations/gabrtv"
	cache_type               'BasicFile'
	cache_options( :path => "#{ENV['HOME']}/.chef/checksums" )
	cookbook_path            ["#{ENV['HOME']}/workspace/deis-chef/cookbooks"]

	knife[:aws_access_key_id] = "#{ENV['AWS_ACCESS_KEY']}"
	knife[:aws_secret_access_key] = "#{ENV['AWS_SECRET_KEY']}"
	knife[:region] = "us-west-2"

Make sure `cookbook_path` points to your checkout of the deis-chef repository.

##### Test Knife connectivity

Use the following command to ensure knife is working properly

    $ knife client list
    gabrtv-validator

##### Cleanup stale instances

If you notice any stale "node" and "client" records, make sure you delete the
instances with `knife ec2 server list` and `knife ec2 server delete`, then run
`bin/cleanup-all.sh` to wipe the Chef "node" and "client" records.

##### Create required key pairs and security groups

* Create an `deis-build` security group with rules
  * 22 open to `0.0.0.0/0` (for SSH admin)
  * 80 open to `0.0.0.0/0` (for slug downloads)
* Create an `deis-instance` security group with rules
  * 22 open to `0.0.0.0/0` (for SSH admin)
  * 80 open to `0.0.0.0/0` (for Nginx proxy access)
  * (optional) 5001-6000 open to `<your-ip>/32` (for app troubleshooting)
* Create an `deis-build` keypair and download the .pem file to `~/Downloads`
* Create an `deis-instance` keypair and download the .pem file to `~/Downloads`

##### Upload latest cookbooks

    $ knife cookbook upload -a
    Uploading apt          [1.9.2]
    Uploading deis         [0.1.0]
    Uploading sudo         [2.1.3]
    Uploaded all cookbooks.

##### Create initial data bags

Create data bags with an empty formation

    $ bin/install-data-bags.sh formation0

##### Provision build system

    $ bin/provision-build.sh

..wait for the Chef bootstrap to finish (takes ~ 10 minutes).

##### Provision new runtime instances

Open one terminal window and run

    $ bin/provision-instance01.sh

Open a second terminal window and run

    $ bin/provision-instance02.sh

..wait for Chef bootstrap to finish on both instances (takes ~ 5 minutes)

## Seed Chef Config

In order to get the server limping along, we need to seed Chef configuration
for each user who will be pushing code to the build server.

Place the temporary Chef config in the user's repo root:

    ssh <build-server>
    sudo -u git -i bash     # get a bash shell as the git user
    mkdir -p /opt/deis/gitosis/repositories/gabrtv/.chef

Write out a `knife.rb` file that knife can use to connect to the Chef server:

    cat > /opt/deis/gitosis/repositories/gabrtv/.chef/knife.rb <<EOF
    current_dir = File.dirname(__FILE__)
    log_level                :info
    log_location             STDOUT
    node_name                "gabrtv"
    client_key               "#{current_dir}/gabrtv.pem"
    chef_server_url          "https://chef.deis.net/organizations/gabrtv"
    EOF

Write out the client key so knife can make API calls to the Chef server:

    cat > /opt/deis/gitosis/repositories/gabrtv/.chef/gabrtv.pem <<EOF
    -----BEGIN RSA PRIVATE KEY-----
    MIIEowIBAAKCAQEApEWdo1wRDmP1W4TLcQbeoiflzsiEhmIsoYCPYQ1GxEMahVtmmX8ElfrvqX78
    . . .
    -----END RSA PRIVATE KEY-----
    EOF

Write out the `deis-instance` SSH key used to trigger Chef converges:

    cat > /opt/deis/gitosis/repositories/gabrtv/.chef/deis-instance.pem <<EOF
    -----BEGIN RSA PRIVATE KEY-----
    MIIEowIBAAKCAQEApEWdo1wRDmP1W4TLcQbeoiflzsiEhmIsoYCPYQ1GxEMahVtmmX8ElfrvqX78
    . . .
    -----END RSA PRIVATE KEY-----
    EOF


## Next Steps


License
-------

**Deis** is open source software under the Apache 2.0 license.
Please see the **LICENSE** file in the root directory for details.


Credits
-------

**Deis** rests on the shoulders of leading open source technologies:

  * Docker
  * Chef
  * Django
  * Heroku buildpacks
  * Gitosis

[OpDemand](http://www.opdemand.com/) sponsors and maintains the
**Deis** project.
