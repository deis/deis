# Deis cluster local test environment

The Vagrantfile and Makefile provided in this directory will launch a multi-host Deis CoreOS cluster that runs locally.

## Prerequisites
On your workstation:
* Install [Vagrant](http://www.vagrantup.com/downloads.html) and [VirtualBox](https://www.virtualbox.org/wiki/Downloads)
* Install [Go](http://golang.org/doc/install) and configure your GOPATH, if necessary
* Install the fleetctl client: `go get github.com/coreos/fleet && go install github.com/coreos/fleet/fleetctl`
* Set the `DEIS_NUM_INSTANCES` environment variable if you'd like more (or less) than the default of 3 machines to test:
  * `export DEIS_NUM_INSTANCES=5`
* Configure the fleetctl client to tunnel through one of the VMs:
  * `export FLEETCTL_TUNNEL=172.17.8.100`
  * (Note that IP addressing for the VMs starts at .100, but you can connect to any VM in the cluster)

## Launching Deis
Follow the normal instructions in the [README](../../README.md) for launching and using Deis, with the caveat that
commands should be run in this directory to ensure that the proper Vagrantfile and Makefile is used. In this directory
pull/build commands are run in each VM, and service starting/stopping are run with `fleetctl` and affect the entire cluster.

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
