Provision a Deis Controller on Rackspace
========================================

1. Install [knife-rackspace][kniferack] with `gem install knife-rackspace` or just `bundle install` from the root directory of your deis repository:

    ```console
    $ cd $HOME/projects/deis
    $ gem install knife-rackspace
    Fetching: knife-rackspace-0.8.1.gem (100%)
    Successfully installed knife-rackspace-0.8.1
    1 gem installed
    Installing ri documentation for knife-rackspace-0.8.1...
    Installing RDoc documentation for knife-rackspace-0.8.1...
    ```

1. Export your Rackspace credentials as environment variables and edit knife.rb to read them:

    ```console
    $ cat <<'EOF' >> $HOME/.bash_profile
    export RACKSPACE_USERNAME=<your_rackspace_username>
    export RACKSPACE_API_KEY=<your_rackspace_api_key>
    EOF
    $ source $HOME/.bash_profile
    $ cat <<'EOF' >> $HOME/.chef/knife.rb
    knife[:rackspace_api_username] = "#{ENV['RACKSPACE_USERNAME']}"
    knife[:rackspace_api_key] = "#{ENV['RACKSPACE_API_KEY']}"
    EOF
    $ knife rackspace server list
    Instance ID  Name  Public IP  Private IP  Flavor  Image  State
    ```

1. Prepare a new server
    1. Create a server named `deis-prepare-image` using the Ubuntu 12.04 LTS image, performance1-2, 1GB performance server
    1. SSH in as root with the password shown
    1. Install the 3.8 kernel with: ```apt-get update && apt-get install -yq linux-image-generic-lts-raring linux-headers-generic-lts-raring && reboot```
    1. After reboot is complete, SSH back in as root and `uname -r` to confirm kernel is `3.8.0-35-generic`
    1. Run the `prepare-controller-image.sh` script to optimize the image for fast boot times

        ```console
        $ \curl -sSL https://raw2.github.com/opdemand/deis/master/contrib/rackspace/prepare-controller-image.sh | bash
        Reading package lists... Done
        Building dependency tree
        Reading state information... Done
        ...
        ```

1. Create a new image from the `deis-prepare-image` server named `deis-node-image`.
    1. In the server list in the Control Panel click the action cog for `deis-prepare-image`
    1. Select "Create New Image" name that image `deis-node-image`
    1. (optionally) Distribute the image to other regions
    1. (optionally) Create/update your Deis flavors to use your new images

1. Make sure to add the 'deis-controller' and the '<your_username>-validator' usernames to the Chef 'admins' group.
    * If you are using hosted Chef, you may need to use the older console to do this: <https://manage.opscode.com/groups/admins/edit>

1. Back on your machine with deis cloned and the deis CLI installed, run the provisioning script to create a new Deis controller:
    * Change ```<region>``` to match the region your image is in (we will add SYD and HKG as soon as performance flavors are available there):
        * dfw
        * ord
        * iad
        * lon

        ```console
        $ cd deis
        $ ./bundle install # if you have not already done so
        $ ./contrib/rackspace/provision-rackspace-controller.sh <region>
        Provisioning a deis controller on Rackspace...
        Creating new SSH key: deis-controller
        + ssh-keygen -f /home/myuser/.ssh/deis-controller -t rsa -N '' -C deis-controller
        + set +x
        Saved to /home/myuser/.ssh/deis-controller
        Created data_bag[deis-formations]
        Created data_bag[deis-apps]
        Provisioning deis-controller with knife rackspace...
        + knife rackspace server create -y --server-create-timeout 1200 --server-name deis-controller --image 4b7c635d-89e1-44be-a15f-2877b5a660d1 --rackspace-region <region> --flavor 4 --identity-file /home/myuser/.ssh/deis-controller --bootstrap-version 11.4.4 --node-name deis-controller --run-list 'recipe[deis::controller]'
        Instance ID: de17ca36-f186-4cdd-8969-4be58e7108ea
        Name: deis-controller
        Flavor: 2GB Standard Instance
        ...
        ```

[kniferack]: http://docs.opscode.com/plugin_knife_rackspace.html
