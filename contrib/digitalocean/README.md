# How to Provision a Deis Controller on Digital Ocean

Here are the steps to get started on Digital Ocean:

## Customize cloud-config.yml
Edit [user-data](../coreos/user-data) and add a discovery URL. This URL will be used by all nodes in this Deis cluster. You can get a new discovery URL by sending a request to http://discovery.etcd.io/new.

## Install DO command line client and authorize:
```console
$ gem install tugboat
$ tugboat authorize
```
You can leave all but the client and API keys as the defaults.

Find out about the ID of your ssh key (import it into DO if it's not listed):
```console
$ tugboat keys
```

## Create a controller image:
```console
$ ./provision-digitalocean-controller-image.sh <YOU SSH KEY ID>
```

## Deploy controllers
Use the created image to launch (an odd number) controller droplets, either via UI
or on the command line using tugboat:
```console
$ tugboat create deis1 -r <REGION_ID> -i <IMAGE_ID> -p true -k <SSH_ID> -s 65
```

Not all regions allow private networks. Choose one which does, e.g. NY 2, Amsterdam 2 or
Singapore 1 at the time of this writing (check the web UI for the current private network
support).

The provisioning script uses a 512 MB droplet by default because for iamge creation
more memory is not needed. Deis controller nodes will need at least 2 GB to even start all
the services. Add the memory requirements of deployed applications and choose an adequat
droplet size, e.g. 8 GB (ID "65").

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
$ export DEIS_HOSTS=10.21.12.1 10.21.12.2 10.21.12.3
```

This variable is used in the `make build` command.
