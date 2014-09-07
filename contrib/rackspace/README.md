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
Copy [user-data.example](../coreos/user-data.example) to [user-data](../coreos/user-data)
and add a new discovery URL.
This URL will be used by all nodes in this Deis cluster. You can get a new discovery URL by sending a request to http://discovery.etcd.io/new.

### Choose number of instances
By default, the provision script will provision 3 servers. You can override this by setting `DEIS_NUM_INSTANCES`:
```console
$ DEIS_NUM_INSTANCES=5 ./provision-rackspace-cluster.sh deis-key
```

Note that for scheduling to work properly, clusters must consist of at least 3 nodes and always have an odd number of members.
For more information, see [optimal etcd cluster size](https://github.com/coreos/etcd/blob/master/Documentation/optimal-cluster-size.md).

Deis clusters of less than 3 nodes are unsupported.

### Run the provision script
Run the [Rackspace provision script](provision-rackspace-cluster.sh) to spawn a new CoreOS cluster.
You'll need to provide the name of the key pair you just added. Optionally, you can also specify a flavor name.
```console
$ cd contrib/rackspace
$ ./provision-rackspace-cluster.sh
Usage: provision-rackspace-cluster.sh <key pair name> [flavor]
$ ./provision-rackspace-cluster.sh deis-key
```

### Choose number of routers
By default, the Makefile will provision 1 router. You can override this by setting `DEIS_NUM_ROUTERS`:
```console
$ export DEIS_NUM_ROUTERS=2
```

### Update CoreOS
Due to image publishing limitations on Rackspace, CoreOS images on Rackspace are frequently out of date.
Each machine needs to be updated. On each one, run:

```console
$ sudo systemctl unmask update-engine && sudo systemctl start update-engine && sudo update_engine_client -update && sudo reboot
```

Once the machine is rebooted, it should have a recent CoreOS version.

### Initialize the cluster
Once the cluster is up, get the hostname of any of the machines from Rackspace, set
FLEETCTL_TUNNEL, and issue a `make run` from the project root:
```console
$ export FLEETCTL_TUNNEL=23.253.219.94
$ cd ../.. && make run
```
The script will deploy Deis and make sure the services start properly.

### Configure DNS
You'll need to configure DNS records so you can access applications hosted on Deis. See [Configuring DNS](http://docs.deis.io/en/latest/installing_deis/configure-dns/) for details.

### Configure Load Balancer
You'll need to create two load balancers on Rackspace to handle your cluster.

    Load Balancer 1
    Port 80
    Protocol HTTP
    Health Monitoring -
      Monitor Type HTTP
      HTTP Path /health-check

    Load Balancer 2
    Virtual IP Shared VIP on Another Load Balancer (select Load Balancer 1)
    Port 2222
    Protocol TCP

### Use Deis!
After that, register with Deis!
```console
$ deis register http://deis.example.org
username: deis
password:
password (confirm):
email: info@opdemand.com
```

## Hack on Deis
If you'd like to use this deployment to build Deis, you'll need to set `DEIS_HOSTS` to an array of your cluster hosts:
```console
$ DEIS_HOSTS="1.2.3.4 2.3.4.5 3.4.5.6" make build
```

This variable is used in the `make build` command.
