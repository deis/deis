# Deis

Deis is an open source PaaS that makes it easy to deploy, scale and manage containers used to host applications and services. Deis builds upon Docker and CoreOS to provide a private PaaS that is lightweight and flexible.

[![Build Status](https://travis-ci.org/deis/deis.png?branch=master)](https://travis-ci.org/deis/deis)
[![Coverage Status](https://coveralls.io/repos/deis/deis/badge.png?branch=master)](https://coveralls.io/r/deis/deis?branch=master)

![Deis Graphic](https://s3-us-west-2.amazonaws.com/deis-images/deis-graphic.png)

# New Deis
Deis has undergone several improvements recently. If you are updating
from Deis version 0.7.0 or earlier, there are several big changes you
should know about. Read the [MIGRATING.md](MIGRATING.md) document for
details.

If you need to use Deis with Chef integration, on Ubuntu 12.04 LTS, or
on DigitalOcean, you should use the
[v0.7.0 release](https://github.com/deis/deis/tree/v0.7.0) of Deis.

# Deploying Deis

Deis is a set of Docker containers that can be deployed anywhere including public cloud, private cloud, bare metal or your workstation. Decide where you'd like to deploy Deis, then follow the deployment-specific documentation for [Rackspace](contrib/rackspace/README.md) or [EC2](contrib/ec2/README.md). Documentation for OpenStack and bare-metal provisioning are forthcoming.

Trying out Deis? Continue following these instructions for a local cluster setup. This is also a great Deis testing/development environment.

## Install prerequisites
On your workstation:
* Install [Vagrant](http://www.vagrantup.com/downloads.html) and [VirtualBox](https://www.virtualbox.org/wiki/Downloads)
* Install the fleetctl client: Install v0.2.0 from the [fleet GitHub page](https://github.com/coreos/fleet/releases).
* Install the Docker client if you want to run Docker commands locally (optional)

## Additional setup for a multi-node cluster
If you'd like to spin up more than one VM to test an entire cluster, there are a few additional prerequisites:
* Set `DEIS_NUM_INSTANCES` to the desired size of your cluster: ```$ export DEIS_NUM_INSTANCES=3```
* Edit [contrib/coreos/user-data](contrib/coreos/user-data) and add a unique discovery URL generated from `https://discovery.etcd.io/new`

## Boot CoreOS

First, start the CoreOS cluster on VirtualBox. From a command prompt, `cd` to the root of the Deis project code and type:

```console
$ vagrant up
```

This instructs Vagrant to spin up each VM. To be able to connect to the VMs, you must add your Vagrant-generated SSH key to the ssh-agent (fleetctl tunnel requires the agent to have this key):
```console
$ ssh-add ~/.vagrant.d/insecure_private_key
```

Export some environment variables so you can connect to the VM using the `docker` and `fleetctl` clients on your workstation.

```console
$ export DOCKER_HOST=tcp://172.17.8.100:4243
$ export FLEETCTL_TUNNEL=172.17.8.100
```

## Build Deis

Use `make pull` to download cached layers from the public Docker Index.  Then use `make build` to assemble all of the Deis components from Dockerfiles.  Grab some coffee while it builds the images on each VM (it can take a while).

```console
$ make pull
$ make build
```

## Run Deis

Use `make run` to start all Deis containers and attach to their log output. This can take some time - the registry service will pull and prepare a Docker image. Grab some more coffee!

```console
$ make run
```

## Additional steps for a multi-node cluster
* Configure local DNS. For a one-node cluster we do this for you: `local.deisapp.com` resolves to the IP of the first VM, 172.17.8.100. Since we cannot know where the `deis-router` container will be running in your cluster, you'll need to setup DNS and resolve a wildcard entry to use for your apps.
* Because of the DNS quandary, we don't start the deis-router component for you. You'll need to start this manually once DNS is setup: `systemctl start deis-router`.

## Install the Deis Client
If you're using the latest Deis release, use `pip install deis` to install the latest [Deis Client](https://pypi.python.org/pypi/deis/) or download [pre-compiled binaries](https://github.com/deis/deis/tree/master/client#get-started).

If you're working off master, precompiled binaries are likely out of date. You should either symlink the python file directly or build a local copy of the client:

```console
$ ln -fs $(pwd)/client/deis.py /usr/local/bin/deis
```
or
```console
$ cd client && python setup.py install
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

# Usage

## Create an Application
Create an application on the default `dev` cluster.

```console
$ deis create
```

Use `deis create --cluster=prod` to place the app on a different cluster.  Don't like our name-generator?  Use `deis create myappname`.

## Push
Push builds of your application from your local git repository or from a Docker Registry.  Each build creates a new release, which can be rolled back.

#### From a Git Repository
When you created the application, a git remote for Deis was added automatically.

```console
$ git push deis master
```
This will use the Deis builder to package your application as a Docker Image and deploy it on your application's cluster.

## Configure
Configure your application with environment variables.  Each config change also creates a new release.

```console
$ deis config:set DATABASE_URL=postgres://
```

Coming soon: Use the integrated ETCD namespace for service discovery between applications on the same cluster.

## Test
### Run tests
Test your application by running commands inside an ephemeral Docker container.

```console
$ deis run make test
```

To integrate with your CI system, check the return code.

### Testing the cluster
Logging into one of the CoreOS machines and stopping a container service should cause the same component on another CoreOS
host to take over as master

These systemd services run the various containers which compose Deis, and can be stopped on a machine with `sudo systemctl stop servicename`.
* deis-builder.service
* deis-cache.service
* deis-controller.service
* deis-database.service
* deis-discovery.service
* deis-logger.service
* deis-registry.service
* deis-router.service

Similarly, bringing down a VM should enable the services on another VM to take over as master

## Scale
Scale containers horizontally with ease.

```console
$ deis scale web=8
```

## Debug
Access to aggregated logs makes it easy to troubleshoot problems with your application.

```console
$ deis logs
```

Use `deis run` to execute one-off commands and explore the deployed container.  Coming soon: `deis attach` to jump into a live container.

## Known Issues

We have sometimes seen the VM reboot while doing `make build` against a
Vagrant virtual machine. If you see this issue using a recent version of
Vagrant and the current master version of Deis, please add to the issue
report at https://github.com/coreos/coreos-vagrant/issues/68 to help us
pin it down.

## License

Copyright 2014, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
