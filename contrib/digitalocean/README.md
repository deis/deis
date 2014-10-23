# Provision a Deis Cluster on DigitalOcean

## Customize cloud-config.yml
Edit [user-data](../coreos/user-data) and add a discovery URL. This URL will be used by all nodes in this Deis cluster. You can get a new discovery URL by sending a request to http://discovery.etcd.io/new.

## Install docl and authorize:
The docl gem consumes the DigitalOcean API.
```console
$ gem install docl
```

Before you can authorize you need to [create a Personal Access Token](https://www.digitalocean.com/community/tutorials/how-to-use-the-digitalocean-api-v2). Make sure you create a read & write token.
Copy paste the token (make sure you also save it somewhere encrypted) before you continue to the next step.

```console
$ docl authorize
```

## Upload keys
Choose an SSH keypair to use for Deis and import it to DigitalOcean:
```console
$ docl upload_key deis ~/.ssh/deis.pub
```
This will print the ID of the uploaded key. Copy and paste the ID, you will need it in a later step.

In case you forget the ID of the public key you can retrieve it later with the following command:
```console
$ docl keys
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
$ ./provision-do-cluster.sh <REGION_ID> <SSH_ID> <SIZE>
```

Not all regions allow private networks. Choose one which does (at the time of this writing, NY 2 & 3,
Amsterdam 2  & 3, Singapore 1 or London 1) - check the web UI for the current private network support.

You can enumerate all of the supported regions with:

```console
$ docl regions --private_networking --metadata
```

Deis controller nodes will need at least 2 GB to even start all
the services. Add the memory requirements of deployed applications and choose an adequate
droplet size. The default is 8 GB. Specify the size with NGB, where N is 2, 4, 8, 16, 32 or 64.

This will print the IP addresses of the initialized machines.

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

## Initialize the cluster
set DEISCTL_TUNNEL to one of the IPs of the virtual machines (these are printed on the console in a previous step). You can also login to the web interface of DigitalOcean to see the Public IP addresses.

```console
$ export DEISCTL_TUNNEL=23.253.219.94
$ deisctl install platform && deisctl start platform
```
Deisctl will deploy Deis and make sure the services are started properly. Grab a coffee.

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
