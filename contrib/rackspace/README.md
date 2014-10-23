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

## Configure Deis
Set the default domain used to anchor your applications:

```console
$ deisctl config platform set domain=mycluster.local
```

For this to work, you'll need to configure DNS records so you can access applications hosted on Deis. See [Configuring DNS](http://docs.deis.io/en/latest/installing_deis/configure-dns/) for details.

If you want to allow `deis run` for one-off admin commands, you must provide an SSH private key that allows Deis to gather container logs on CoreOS hosts:

```console
$ deisctl config platform set sshPrivateKey=<path-to-private-key>
```

### Initialize the cluster
Once the cluster is up, get the hostname of any of the machines from Rackspace, set
DEISCTL_TUNNEL and install the platform:
```console
$ export DEISCTL_TUNNEL=23.253.219.94
$ deisctl install platform && deisctl start platform
```

The installer will deploy Deis and make sure the services start properly.

### Choose number of routers
By default, `deisctl` will provision 1 router. You can override this by scaling up:
```console
$ deisctl scale router=2
```

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

See [Hacking on Deis](http://docs.deis.io/en/latest/contributing/hacking/).
