# Provision a Deis Cluster on DigitalOcean

## Note on Deis support for DigitalOcean
DigitalOcean support was contributed by [sttts](https://github.com/sttts). The CoreOS bootstrapping
is heavily based on [Levi Aul's code](https://gist.github.com/tsutsu/490f35f48897df0f5173).

Note that DigitalOcean does not support CoreOS images natively. To work around this, the provision
scripts for DigitalOcean first create a CoreOS image which will be used for provisioning the cluster.

Until native DigitalOcean support for CoreOS is implemented, this workaround is likely to be more fragile than deploying
Deis to other cloud platforms which support CoreOS natively.

UPDATE: As of July 15, 2014, native CoreOS support on DigitalOcean is [planned](http://digitalocean.uservoice.com/forums/136585-digital-ocean/suggestions/4250154-suport-coreos-as-a-deployment-platform).

## Customize cloud-config.yml
Edit [user-data](../coreos/user-data) and add a discovery URL. This URL will be used by all nodes in this Deis cluster. You can get a new discovery URL by sending a request to http://discovery.etcd.io/new.

## Install tugboat and authorize:
The tugboat gem consumes the DigitalOcean API.
```console
$ gem install tugboat
$ tugboat authorize
```
You can leave all but the client and API keys as the defaults.

## Upload keys
Choose an SSH keypair to use for Deis and import it to DigitalOcean:
```console
$ tugboat add-key deis
```

Then, get the ID of the key:
```console
$ tugboat keys
```

## Create a Deis image:
```console
$ ./provision-digitalocean-deis-image.sh <SSH KEY ID>
```

## Choose number of instances
By default, the script will provision 3 servers. You can override this by setting `DEIS_NUM_INSTANCES`:
```console
$ export DEIS_NUM_INSTANCES=5
```

Note that for scheduling to work properly, clusters must consist of at least 3 nodes and always have an odd number of members.
For more information, see [optimal etcd cluster size](https://github.com/coreos/etcd/blob/master/Documentation/optimal-cluster-size.md).

Deis clusters of less than 3 nodes are unsupported.

## Deploy cluster
Run the provision script:
```console
$ ./provision-do-cluster.sh <REGION_ID> <IMAGE_ID> <SSH_ID> <SIZE>
```

Not all regions allow private networks. Choose one which does (at the time of this writing, NY 2,
Amsterdam 2, Singapore 1 or London 1) - check the web UI for the current private network support.

You can enumerate all the regions with:

```console
$ tugboat regions
```

The provisioning script uses a 512 MB droplet by default because for image creation
more memory is not needed. Deis controller nodes will need at least 2 GB to even start all
the services. Add the memory requirements of deployed applications and choose an adequate
droplet size. The default is 8 GB (ID "65"). You can enumerate all sizes with:

```console
$ tugboat sizes
```

## Choose number of routers
By default, the Makefile will provision 1 router. You can override this by setting `DEIS_NUM_ROUTERS`:
```console
$ export DEIS_NUM_ROUTERS=2
```

## Initialize the cluster
Once the cluster is up, get the IPs of any of the machines using `tugboat droplets`, set
FLEETCTL_TUNNEL to one of these IPs:
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
$ DEIS_HOSTS="1.2.3.4 2.3.4.5 3.4.5.6" make build
```

This variable is used in the `make build` command.
