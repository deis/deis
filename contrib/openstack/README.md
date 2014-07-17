# Provision a Deis Cluster on Openstack

We'll mostly be following the [CoreOS on Openstack](https://coreos.com/docs/running-coreos/cloud-providers/openstack/) guide. You'll need to have a sane python environment with pip already installed (`sudo easy_install pip`).

### Install supernova and its dependencies:
```console
ubuntu
$ sudo apt-get install python-novaclient
or other system
$ sudo pip install https://github.com/openstack/python-novaclient/archive/2.18.1.zip
```

### Configure supernova
Edit `~novaclient-env.sh` to match the following:
```
export OS_AUTH_URL = {openstack_keystone}
export OS_USERNAME = {openstack_username}
export OS_PASSWORD = {openstack_api_key}
export OS_TENANT_NAME = {openstack_account_id}
```

Your account ID is displayed in the upper right-hand corner of the cloud control panel UI, and your API key can be found on the Account Settings page.

### Set up your keys
Choose an existing keypair or generate a new one, if desired. Tell supernova about the key pair and give it an identifiable name:

```console
$ nova keypair-add --pub-key ~/.ssh/deis.pub deis-key
```

### Customize cloud-config.yml
Edit [user-data](../coreos/user-data) and add a discovery URL. This URL will be used by all nodes in this Deis cluster. You can get a new discovery URL by sending a request to http://discovery.etcd.io/new.

### Choose number of instances
By default, the provision script will provision 3 servers. You can override this by setting `DEIS_NUM_INSTANCES`:
```console
$ export DEIS_NUM_INSTANCES=5
```

Note that for scheduling to work properly, clusters must consist of at least 3 nodes and always have an odd number of members.
For more information, see [optimal etcd cluster size](https://github.com/coreos/etcd/blob/master/Documentation/optimal-cluster-size.md).

Deis clusters of less than 3 nodes are unsupported.

### Run the provision script
Run the [Openstack provision script](provision-openstack-cluster.sh) to spawn a new CoreOS cluster.
You'll need to provide the name of the key pair you just added. Optionally, you can also specify a flavor name.

```console
$ cd contrib/openstack
Choose: openstack image id
$ nova image-list
Choose: openstack flavor id
$ nova flavor-list
Choose: openstack net id
$ nova net-list
Usage: provision-openstack-cluster.sh <image-id> <flavor-id> deis-key <net-id>
$ ./provision-openstack-cluster.sh  <image-id> <flavor-id> deis-key <net-id>
```

### Choose number of routers
By default, the Makefile will provision 1 router. You can override this by setting `DEIS_NUM_ROUTERS`:
```console
$ export DEIS_NUM_ROUTERS=2
```

### Initialize the cluster
Once the cluster is up, get the hostname of any of the machines from Openstack, set
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
$ export DEIS_HOSTS=10.21.12.1 10.21.12.2 10.21.12.3
```

This variable is used in the `make build` command.
