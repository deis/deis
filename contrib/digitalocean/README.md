How to Provision a Deis Controller on Digital Ocean
===================================================

Here are the steps to get started on Digital Ocean:

* install [knife-digital_ocean][kdo]

```
bundle install
```

* install python requirements

```
pip install -r requirements.txt
```

* add the following to ~/.chef/knife.rb

```
knife[:digital_ocean_client_id] =   "your digital ocean client ID"
knife[:digital_ocean_api_key] =     "your digital ocean API key"
```

* Follow the steps provided in contrib/digitalocean/prepare-digitalocean-snapshot.sh
* Run this command to start the provisioning process

```
./contrib/digitalocean/provision-digitalocean-controller.sh
```

This script will read from your knife config file, create a new SSH key for the controller, upload the SSH key to digital ocean, and provision a Deis Controller from the snapshot created.

[kdo]: https://github.com/rmoriz/knife-digital_ocean
