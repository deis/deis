# Provision a Deis Cluster on OpenStack

**NOTE**
OpenStack support for Deis was contributed by @shakim. OpenStack support is untested by the Deis team, so we rely on the community to improve these scripts and to fix bugs.
We greatly appreciate the help!

### Prerequisites:
Make sure that the following utilities are installed and in your execution path:
- nova
- neutron
- glance

### Install Deis CLI tools

```console
$ sudo pip install deis
$ curl -sSL http://deis.io/deisctl/install.sh | sh -s 1.0.1
$ mv deisctl /usr/local/bin
$ chmod +x /usr/local/bin/deisctl
```

### Configure openstack
Create an `openrc.sh` file to match the following:
```
export OS_AUTH_URL={openstack_auth_url}
export OS_USERNAME={openstack_username}
export OS_PASSWORD={openstack_password}
export OS_TENANT_NAME={openstack_tenant_name}
```

(Alternatively, download OpenStack RC file from Horizon/Access & Security/API Access.)

Source your nova credentials:

```console
$ source openrc.sh
```

### Set up your keys
Choose an existing keypair or upload a new public key, if desired.

```console
$ nova keypair-add --pub-key ~/.ssh/deis.pub deis-key
```

### Upload a coreos image to Glance

You need to have a relatively recent CoreOS image.  If you don't have one and your Openstack install allows you to upload your own images you can do the following:

```console
$ wget http://alpha.release.core-os.net/amd64-usr/current/coreos_production_openstack_image.img.bz2
$ bunzip2 coreos_production_openstack_image.img.bz2
$ glance image-create --name coreos \
  --container-format bare \
  --disk-format qcow2 \
  --file coreos_production_openstack_image.img \
  --is-public True
```

### Customize user-data

Create a user-data file with a new discovery URL this way:

```console
$ make discovery-url
```

Or copy [`contrib/coreos/user-data.example`](../coreos/user-data.example) to `contrib/coreos/user-data` and follow the directions in the `etcd:` section to add a unique discovery URL.

### Choose number of instances and routers

```console
$ export DEIS_NUM_INSTANCES=3
$ export DEIS_NUM_ROUTERS=1
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

If you have a fairly straight forward openstack install you should be able to use the provisioning script provided.   This script assumes you are using neutron and have security-groups enabled.

Run the [Openstack provision script](provision-openstack-cluster.sh) to spawn a new CoreOS cluster.
You'll need to provide the name of the CoreOS image name (or ID), and the key pair you just added. Optionally, you can also specify a flavor name.
```console
$ cd contrib/openstack
$ ./provision-openstack-cluster.sh
Usage: provision-openstack-cluster.sh <coreos image name/id> <key pair name> [flavor]
$ ./provision-openstack-cluster.sh coreos deis-key
```

You can override the name of the internal network to use by setting the environment variable `DEIS_NETWORK=internal`.  If this doesn't exist the script will try to create it with the default CIDR which requires your openstack cluster to support tenant vlans.

You can also override the name of the security group to attach to the instances by setting `DEIS_SECGROUP=deis_test`.  If this doesn't exist the script will attempt to create it.  If you are creating your own security groups you can use the provision script as a guide.  Make sure that you have a rule to enable full communication inside the security group, or you will have a bad day.

### Manually start the instances

### Finish of your openstack configuration by setting up floating IPs.

You will want to attach a floating ip to at least one of your instances.  You'll do that like this:

```
$ nova floating-ip-create <pool>
$  nova floating-ip-associate deis-1 <IP provided by above command>
```

### Initialize the cluster
Once the cluster is up:
* **If required, allocate and associate floating IPs to any or all of your hosts**
* Get the IP address of any of the machines from Openstack
* Set the default domain used to anchor your applications:

```console
$ deisctl config platform set domain=mycluster.local
```

** For this to work, you'll need to configure DNS records so you can access applications hosted on Deis. See [Configuring DNS](http://docs.deis.io/en/latest/managing_deis/configure-dns/#dns-records) for details.

* If you want to allow `deis run` for one-off admin commands, you must provide an SSH private key that allows Deis to gather container logs on CoreOS hosts:

```console
$ deisctl config platform set sshPrivateKey=<path-to-private-key>
```

* set DEISCTL_TUNNEL to one of your floating IPs and install the platform:

```console
$ export DEISCTL_TUNNEL=<Floating IP>
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
