# Provision a Deis Cluster on OpenStack

**NOTE**
OpenStack support for Deis was contributed by @shakim. OpenStack support is untested by the Deis team, so we rely on the community to improve these scripts and to fix bugs.
We greatly appreciate the help!

### Prerequisites:
Make sure that the following utilities are installed and in your execution path:
- nova
- neutron

### Configure nova
Create an `openrc.sh` file to match the following:
```
[production]
OS_AUTH_URL = {openstack_auth_url}
OS_USERNAME = {openstack_username}
OS_PASSWORD = {openstack_api_key}
OS_TENANT_ID = {openstack_tenant_id}
OS_TENANT_NAME = {openstack_tenant_name}
```

(Alternatively, download OpenStack RC file from Horizon/Access & Security/API Access.)

Source your nova credentials:

```console
# source openrc.sh
```

### Set up your keys
Choose an existing keypair or upload a new public key, if desired.

```console
$ nova keypair-add --pub-key ~/.ssh/deis.pub deis-key
```

### Customize cloud-config.yml
Edit [user-data](../coreos/user-data) and add a discovery URL. This URL will be used by all nodes in this Deis cluster. You can get a new discovery URL by sending a request to http://discovery.etcd.io/new.

### Choose number of instances
By default, the provision script will provision 3 servers. You can override this by setting `DEIS_NUM_INSTANCES`:
```console
$ DEIS_NUM_INSTANCES=5 ./provision-openstack-cluster.sh deis-key
```

Note that for scheduling to work properly, clusters must consist of at least 3 nodes and always have an odd number of members.
For more information, see [optimal etcd cluster size](https://github.com/coreos/etcd/blob/master/Documentation/optimal-cluster-size.md).

Deis clusters of less than 3 nodes are unsupported.

### Deis network settings
The script creates a private network called 'deis' if no such network exists.

By default, the deis subnet IP range is set to 10.21.12.0/24. To override it and the default DNS settings, set the following variables:
```console
$ export DEIS_CIDR=10.21.12.0/24
$ export DEIS_DNS=10.21.12.3,8.8.8.8
```

**_Please note that this script does not handle floating IPs or routers. These should be provisioned manually either by Horizon or CLI_**

### Run the provision script
Run the [Openstack provision script](provision-openstack-cluster.sh) to spawn a new CoreOS cluster.
You'll need to provide the name of the CoreOS image name (or ID), and the key pair you just added. Optionally, you can also specify a flavor name.
```console
$ cd contrib/openstack
$ ./provision-openstack-cluster.sh
Usage: provision-openstack-cluster.sh <coreos image name/id> <key pair name> [flavor]
$ ./provision-openstack-cluster.sh coreos deis-key
```

### Choose number of routers
By default, the Makefile will provision 1 router. You can override this by setting `DEIS_NUM_ROUTERS`:
```console
$ export DEIS_NUM_ROUTERS=2
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
Once the cluster is up:
* **If required, allocate and associate floating IPs to any or all of your hosts**
* Get the IP address of any of the machines from Openstack
* set DEISCTL_TUNNEL and install the platform:

```console
$ export DEISCTL_TUNNEL=23.253.219.94
$ deisctl install platform && deisctl start platform
```

The installer will deploy Deis and make sure the services start properly.

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
