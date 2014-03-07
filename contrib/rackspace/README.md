Provision a Deis Controller on Rackspace
========================================

1. Install [knife-rackspace][kniferack] with `gem install knife-rackspace` or just `bundle install` from the root directory of your deis repository:

    ```console
    $ cd $HOME/projects/deis
    $ gem install knife-rackspace
    Fetching: knife-rackspace-0.9.0.gem (100%)
    Successfully installed knife-rackspace-0.9.0
    1 gem installed
    Installing ri documentation for knife-rackspace-0.9.0...
    Installing RDoc documentation for knife-rackspace-0.9.0...
    ```

2. Export your Rackspace credentials as environment variables and edit knife.rb to read them:

    ```console
    $ cat <<'EOF' >> $HOME/.bash_profile
    export RACKSPACE_USERNAME=<your_rackspace_username>
    export RACKSPACE_API_KEY=<your_rackspace_api_key>
    EOF
    $ cat <<'EOF' > $HOME/.rackspacerc
    export OS_AUTH_URL="https://identity.api.rackspacecloud.com/v2.0/"
    export OS_USERNAME=$RACKSPACE_USERNAME
    export OS_PASSWORD=$RACKSPACE_API_KEY
    export OS_TENANT_NAME=$RACKSPACE_USERNAME
    export OS_TENANT_ID=<your_rackspace_tenant_id>
    export OS_REGION_NAME=<your_rackspace_region_name>
    export OS_AUTH_SYSTEM="rackspace"
    EOF
    $ source $HOME/.bash_profile
    $ source $HOME/.rackspacerc
    $ cat <<'EOF' >> $HOME/.chef/knife.rb
    knife[:rackspace_api_username] = "#{ENV['RACKSPACE_USERNAME']}"
    knife[:rackspace_api_key] = "#{ENV['RACKSPACE_API_KEY']}"
    knife[:rackspace_ssh_keypair]   = "deis"
    knife[:rackspace_region]        = #{ENV['OS_REGION_NAME']}
    EOF
    $ knife rackspace server list
    Instance ID  Name  Public IP  Private IP  Flavor  Image  State
    ```

3. Now you can follow the standard deis setup:
  ```bash
  bundle install # Installs gem files like the knife tool
  berks install # Downloads the relevant cookbooks
  # '--ssl-verify' is only needed when using a self-hosted Chef Server
  # hint: you can also set that at $HOME/.berkshelf/config.json
  berks upload [--ssl-verify=false] # Upload the cookbooks to the Chef Server
  ```

4. Prepare a new server
    1. Create a server named `deis-prepare-image` using the Ubuntu 12.04 LTS image, performance1-2, 1GB performance server
    2. SSH in as root with the password shown
    3. Install the 3.11 kernel with: ```apt-get update && apt-get install -yq linux-image-generic-lts-saucy linux-headers-generic-lts-saucy && reboot```
    4. After reboot is complete, SSH back in as root and `uname -r` to confirm kernel is `3.11.0-17-generic`
    5. Run the `prepare-controller-image.sh` script to optimize the image for fast boot times

        ```console
        $ ssh root@ip-address 'bash -s' < contrib/rackspace/prepare-node-image.sh
        Reading package lists... Done
        Building dependency tree
        Reading state information... Done
        ...
        ```

5. Create a new image from the `deis-prepare-image` server named `deis-node-image`.
    1. In the server list in the Control Panel click the action cog for `deis-prepare-image`
    2. Select "Create New Image" name that image `deis-node-image`
    3. (optionally) Distribute the image to other regions
    4. (optionally) Create/update your Deis flavors to use your new images

6. Make sure to add the `deis-controller` and the `<your_username>-validator` usernames to the Chef 'admins' group.
    * If you are using hosted Chef, you may need to use the older console to do this: <https://manage.opscode.com/groups/admins/edit>

7. Back on your machine with deis cloned and the deis CLI installed, run the provisioning script to create a new Deis controller:
    * Change ```<region>``` to match the region your image is in (we will add SYD and HKG as soon as performance flavors are available there):
        * dfw
        * ord
        * iad
        * lon

        ```console
        $ cd deis
        $ bundle install # if you have not already done so
        $ ./contrib/rackspace/provision-rackspace-controller.sh <region>
        Provisioning a deis controller on Rackspace...
        Creating new SSH key: id_rsa
        + ssh-keygen -f /home/deis/.ssh/id_rsa -t rsa -N '' -C deis
        + set +x
        Saved to /home/deis/.ssh/id_rsa
        Created data_bag[deis-formations]
        Created data_bag[deis-apps]
        Provisioning deis-controller-H7WVl with knife rackspace...
        + knife rackspace server create --bootstrap-version 11.8.2 --rackspace-region ord --image f569b831-afe5-44f5-85eb-3bf9e1d0d336 --flavor performance1-2 --rackspace-metadata '{"Name": "deis-controller-H7WVl"}' --rackspace-disk-config MANUAL --server-name deis-controller-H7WVl --node-name deis-controller-H7WVl --run-list 'recipe[deis::controller]'
        Instance ID: cf7aeadd-4bb1-4f69-9238-7a0586a863b9
        Name: deis-controller-H7WVl
        Flavor: 2 GB Performance
        Image: deis-node-image
        Metadata: [  <Fog::Compute::RackspaceV2::Metadatum
            key="Name",
            value="deis-controller-H7WVl"
          >]
        RackConnect Wait: no
        ServiceLevel Wait: no
        SSH Key: deis
        ...
        ```

[kniferack]: http://docs.opscode.com/plugin_knife_rackspace.html
