# Deis local single-host cluster

The Vagrantfile and Makefile provided in this directory will launch a single-host Deis CoreOS cluster that runs locally. To test a multiple-node cluster, see the [multi-host Vagrant guide](../vagrant-n-node/README.md).

All commands below are assumed to be run in this directory.

## Prerequisites
On your workstation:
* Install [Vagrant](http://www.vagrantup.com/downloads.html) and [VirtualBox](https://www.virtualbox.org/wiki/Downloads)

## Boot CoreOS

First, start a CoreOS virtual machine on VirtualBox.

```console
$ vagrant up
```

## Build Deis

Use `make pull` to download cached layers from the public Docker Index.  Then use `make build` to assemble all of the Deis components from Dockerfiles.  Grab some coffee while it builds the images on the CoreOS VM (it can take a while).

```console
$ make pull
$ make build
```

## Run Deis

Use `make run` to start all Deis containers and attach to their log output.

```console
$ make run
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
