# Deis local multi-host cluster

The Vagrantfile and Makefile provided in this directory will launch a multi-host Deis CoreOS cluster that runs locally. If you just want to get Deis up and running quickly, a [single-host](../vagrant-1-node/README.md) guide is also provided.

The multi-host Vagrant setup differs from the single-host setup in that it creates a resilient, multi-host cluster that closely resembles how Deis would behave in a production environment.

## Prerequisites
On your workstation:
* Install [Vagrant](http://www.vagrantup.com/downloads.html) and [VirtualBox](https://www.virtualbox.org/wiki/Downloads)
* Install [Go](http://golang.org/doc/install) and configure your GOPATH, if necessary
* Install the fleetctl client: `go get github.com/coreos/fleet && go install github.com/coreos/fleet/fleetctl`
* Set the `DEIS_NUM_INSTANCES` environment variable if you'd like more (or less) than the default of 3 machines to test:
  * `export DEIS_NUM_INSTANCES=5`

## Boot the VMs

The Vagrantfile will define `DEIS_NUM_INSTANCES` VMs. To start them all, simply type:

```console
$ vagrant up
```

Export some environment variables so you can connect to the VM using the `docker` and [`fleetctl`](https://github.com/coreos/fleet#building) clients on your workstation.

```console
$ export DOCKER_HOST=tcp://172.17.8.100:4243
$ export FLEETCTL_TUNNEL=172.17.8.100
```

(Note that IP addressing for the VMs starts at .100, and you can connect to any VM in the cluster)

## Build Deis

Use `make pull` to download cached layers from the public Docker Index.  Then use `make build` to assemble all of the Deis components from Dockerfiles.  Grab some coffee while it builds the images on each of the CoreOS VM (it can take a while).

```console
$ make pull
$ make build
```

## Run Deis

Use `make run` to start all Deis containers. This causes fleet to schedule each job on one of the VMs:

```console
$ make run
Job deis-registry.service scheduled to 73c7d285.../172.17.8.100
Job deis-logger.service scheduled to 21ad134c.../172.17.8.101
Job deis-database.service scheduled to 73c7d285.../172.17.8.100
Job deis-cache.service scheduled to 73c7d285.../172.17.8.100
Job deis-controller.service scheduled to e5c14be6.../172.17.8.102
Job deis-builder.service scheduled to e5c14be6.../172.17.8.102
Job deis-router.service scheduled to 73c7d285.../172.17.8.100
```

## Register a User

Use the Deis Client to register a new user.

```console
$ deis register http://local.deisapp.com:8000
$ deis keys:add
```

Use `deis keys:add` to add your SSH public key for `git push` access.

## Initalize a Cluster

Initalize a `dev` cluster with a list of CoreOS hosts and your CoreOS private key.

```console
$ deis clusters:create dev local.deisapp.com --hosts=local.deisapp.com --auth=~/.vagrant.d/insecure_private_key
```

The `dev` cluster will be used as the default cluster for future `deis` commands.

## Use Deis!

Follow the instructions in the [README](../../README.md) to deploy your first application, or check out the [docs](http://docs.deis.io/en/latest/) for more information. Have fun!

## Testing ideas
### Stop a container service
Logging into one of the CoreOS machines and stopping a container service should cause the same component on another CoreOS
host to take over as master

### Stop a VM
Similarly, bringing down a VM should enable the services on another VM to take over as master

## Useful references
### systemd services
These systemd services run the various containers which compose Deis, and can be stopped on a machine with `sudo systemctl stop servicename`.
* deis-builder.service
* deis-cache.service
* deis-controller.service
* deis-database.service
* deis-discovery.service
* deis-logger.service
* deis-registry.service
* deis-router.service
