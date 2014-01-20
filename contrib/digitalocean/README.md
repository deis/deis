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

* Follow theses steps to create a deis-base snapshot.

[install Packer](http://www.packer.io/intro/getting-started/setup.html).
Then, `cd` into ./contrib/digitalocean/ and run the following:


```
$ export DIGITALOCEAN_CLIENT_ID="YOUR CLIENT ID"
$ export DIGITALOCEAN_API_KEY="YOUR API KEY"
$ packer build deis-base.json
```


* Run this command to start the provisioning process

```
./contrib/digitalocean/provision-digitalocean-controller.sh
```

This script will read from your knife config file, create a new SSH key for the controller, upload the SSH key to digital ocean, and provision a Deis Controller from the snapshot created.

[kdo]: https://github.com/rmoriz/knife-digital_ocean
