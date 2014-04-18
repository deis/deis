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
$ ./provision-rackspace-cluster.sh
Usage: provision-rackspace-cluster.sh <key pair name> [flavor]
$ ./provision-rackspace-cluster.sh deis-key
```

By default, the script will provision 3 servers. You can override this by setting `DEIS_NUM_INSTANCES`:
```console
$ DEIS_NUM_INSTANCES=5 ./provision-rackspace-cluster.sh deis-key
```

### Initialize the cluster
Once the cluster is up, get the IP address for any of the machines in the cluster, set
FLEETCTL_TUNNEL, and run [the init script](initialize-rackspace-cluster.sh) to bootstrap the cluster
remotely:
```console
$ export FLEETCTL_TUNNEL=23.253.219.94
$ ./initialize-rackspace-cluster.sh
The authenticity of host '23.253.219.94:22' can't be established.
RSA key fingerprint is ce:3a:c1:3a:ad:11:bd:60:84:8e:60:a8:2f:19:1a:a6.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '23.253.219.94:22' (RSA) to the list of known hosts.
Job deis-registry.service scheduled to 73c7d285.../23.253.218.114
Job deis-logger.service scheduled to 21ad134c.../23.253.217.229
Job deis-database.service scheduled to 73c7d285.../23.253.218.114
Job deis-cache.service scheduled to 73c7d285.../23.253.218.114
Job deis-controller.service scheduled to e5c14be6.../23.253.219.94
Job deis-builder.service scheduled to e5c14be6.../23.253.219.94
Job deis-router.service scheduled to 73c7d285.../23.253.218.114
done!
```

### Use Deis!
After that, wait for the components to come up, check which host the controller is
running on and register with Deis!
```
$ fleetctl list-units
UNIT                    LOAD    ACTIVE  SUB     DESC            MACHINE
deis-builder.service    loaded  active  running deis-builder    e5c14be6.../23.253.219.94
deis-cache.service      loaded  active  running deis-cache      73c7d285.../23.253.218.114
deis-controller.service loaded  active  running deis-controller e5c14be6.../23.253.219.94
deis-database.service   loaded  active  running deis-database   73c7d285.../23.253.218.114
deis-logger.service     loaded  active  running deis-logger     21ad134c.../23.253.217.229
deis-registry.service   loaded  active  running deis-registry   73c7d285.../23.253.218.114
deis-router.service     loaded  active  running deis-router     73c7d285.../23.253.218.114

$ deis register 23.253.219.94:8000
username: deis
password:
password (confirm):
email: info@opdemand.com
```
