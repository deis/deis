# Provision a Deis Cluster on Rackspace

We'll mostly be following the [CoreOS on Rackspace](https://coreos.com/docs/running-coreos/cloud-providers/rackspace/) guide. You'll need to have a sane python environment with pip already installed (`sudo easy_install pip`).

### Install supernova and its dependencies:
```console
$ sudo pip install keyring
$ sudo pip install rackspace-novaclient
$ sudo pip install supernova
```

### Configure supernova
Edit `~/.supernova` to match the following:
```
[production]
OS_AUTH_URL = https://identity.api.rackspacecloud.com/v2.0/
OS_USERNAME = {rackspace_username}
OS_PASSWORD = {rackspace_api_key}
OS_TENANT_NAME = {rackspace_account_id}
OS_REGION_NAME = DFW (or ORD or another region)
OS_AUTH_SYSTEM = rackspace
```

Your account ID is displayed in the upper right-hand corner of the cloud control panel UI, and your API key can be found on the Account Settings page.

### Set up your keys
Choose an existing keypair or generate a new one, if desired. Tell supernova about the key pair and give it an identifiable name:

```console
$ supernova production keypair-add --pub-key ~/.ssh/deis.pub deis-key
```

### Customize cloud-config.yml
Edit [user-data](../coreos/user-data) and add a discovery URL. This URL will be used by all nodes in this Deis cluster. You can get a new discovery URL by sending a request to http://discovery.etcd.io/new.

### Run the provision script
Run the [Rackspace provision script](provision-rackspace-cluster.sh) to spawn a new CoreOS cluster.
You'll need to provide the name of the key pair you just added. Optionally, you can also specify a flavor name.
```console
$ cd contrib/rackspace
$ ./provision-rackspace-cluster.sh
Usage: provision-rackspace-cluster.sh <key pair name> [flavor]
$ ./provision-rackspace-cluster.sh deis-key
```

By default, the script will provision 3 servers. You can override this by setting `DEIS_NUM_INSTANCES`:
```console
$ DEIS_NUM_INSTANCES=5 ./provision-rackspace-cluster.sh deis-key
```

### Initialize the cluster
Once the cluster is up, get the hostname of any of the machines from EC2, set
FLEETCTL_TUNNEL, and issue a `make run` from the project root:
```console
$ export FLEETCTL_TUNNEL=23.253.219.94
$ cd ../.. && make run
```
The script will deploy Deis and make sure the services start properly.

### Configure DNS
You'll need to configure DNS records so you can access applications hosted on Deis. See [Configuring DNS](http://docs.deis.io/en/latest/operations/configure-dns/) for details.

### Use Deis!
After that, register with Deis!
```
$ deis register deis.example.org:8000
username: deis
password:
password (confirm):
email: info@opdemand.com
```
