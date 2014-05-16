# Deis

Deis (pronounced DAY-iss) is an open source PaaS that makes it easy to deploy and manage applications on your own servers. Deis builds upon [Docker](http://docker.io/) and [CoreOS](http://coreos.com) to provide a lightweight PaaS with a [Heroku-inspired](http://heroku.com) workflow.

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

Deis is a set of Docker containers that can be deployed anywhere including public cloud, private cloud, bare metal or your workstation. Decide where you'd like to deploy Deis, then follow the deployment-specific documentation for [Rackspace](contrib/rackspace/README.md), [EC2](contrib/ec2/README.md), or [bare-metal](contrib/bare-metal/README.md) provisioning. Documentation for other platforms is forthcoming. Want to see a particular platform supported? Open an [issue](https://github.com/deis/deis/issues/new) and we'll investigate.

Trying out Deis? Continue following these instructions for a local cluster setup. This is also a great Deis testing/development environment.

## Install prerequisites
On your workstation:
* Install [Vagrant](http://www.vagrantup.com/downloads.html) and [VirtualBox](https://www.virtualbox.org/wiki/Downloads)
* Install the fleetctl client: Install v0.3.2 from the [fleet GitHub page](https://github.com/coreos/fleet/releases).
* Install the Docker client if you want to run Docker commands locally (optional)

If you're on Ubuntu, you need to install the nfs-kernel-server package as it's required for sharing folders between your host and your CoreOS VM:

    sudo apt-get install nfs-kernel-server

## Additional setup for a multi-node cluster
If you'd like to spin up more than one VM to test an entire cluster, there are a few additional prerequisites:
* Edit [contrib/coreos/user-data](contrib/coreos/user-data) and add a unique discovery URL generated from `https://discovery.etcd.io/new`
* Set `DEIS_NUM_INSTANCES` to the desired size of your cluster (typically 3 or 5): ```$ export DEIS_NUM_INSTANCES=3```
* Instead of `local.deisapp.com`, use either `local3.deisapp.com` or `local5.deisapp.com` as your hostname for accessing the cluster

Note that for scheduling to work properly, clusters must consist of at least 3 nodes and always have an odd number of members.
For more information, see [optimal etcd cluster size](https://github.com/coreos/etcd/blob/master/Documentation/optimal-cluster-size.md).

Deis clusters of less than 3 nodes are unsupported for anything other than local development.

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

Your Vagrant VM is accessible at `local.deisapp.com` (or `local3.deisapp.com`/`local5.deisapp.com`). For clusters on other platforms (EC2, Rackspace, bare metal, etc.), see our guide to [Configuring DNS](http://docs.deis.io/en/latest/operations/configure-dns/).

## Testing the cluster
Integration tests and corresponding documentation can be found under the `test/` folder.

These systemd services run the various containers which compose Deis, and can be stopped on a machine with `sudo systemctl stop servicename`.
* deis-builder.service
* deis-cache.service
* deis-controller.service
* deis-database.service
* deis-discovery.service
* deis-logger.service
* deis-registry.service
* deis-router.service

Logging into one of the CoreOS machines and stopping a container service should cause the same component on another CoreOS
host to take over as master.
Similarly, bringing down a VM should enable the services on another VM to take over as master.

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

## Initialize a Cluster

Initialize a `dev` cluster with a list of CoreOS hosts and your CoreOS private key.

```console
$ deis clusters:create dev local.deisapp.com --hosts=local.deisapp.com --auth=~/.vagrant.d/insecure_private_key
```

The parameters to `deis clusters:create` are:
* cluster name (`dev`) - the name used by Deis to reference the cluster
* cluster hostname (`local.deisapp.com`) - the hostname under which apps are created, like `balancing-giraffe.local.deisapp.com`
* cluster members (`--hosts`) - a comma-separated list of cluster members -- not necessarily all members, but at least one (for EC2 and Rackspace, this is a list of the internal IPs like `--hosts=10.21.12.1,10.21.12.2,10.21.12.3`)
* auth SSH key (`--auth`) - the SSH private key used to provision servers (for EC2 and Rackspace, this key is likely `~/.ssh/deis`)

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

## Troubleshooting

Common issues that users have run into when provisioning Deis are detailed below.

#### When running a `make` action - 'Failed initializing SSH client: ssh: handshake failed: ssh: unable to authenticate'
Did you remember to add your SSH key to the ssh-agent? `ssh-agent -L` should list the key you used to provision the servers. If it's not there, `ssh-add -K /path/to/your/key`.

#### When running a `make` action - 'All the given peers are not reachable (Tried to connect to each peer twice and failed)'
The most common cause of this issue is that a [new discovery URL](https://discovery.etcd.io/new) wasn't generated and updated in [contrib/coreos/user-data](contrib/coreos/user-data) before the cluster was launched. Each Deis cluster must have a unique discovery URL, else there will be entries for old hosts that etcd will try and fail to connect to. Destroy and relaunch the cluster, ensuring to use a fresh discovery URL.

#### Various NFS issues, specifically an 'access denied' error
The easiest workaround for this is to use rsync instead of NFS. In the Vagrantfile, swap out the `config.vm.synced_folder` line for the commented version.

#### Scaling an app doesn't work, and/or the app shows 'Welcome to nginx!'
This means the controller failed to submit jobs for the app to fleet. `fleetctl status deis-controller` will show detailed error information, but the most common cause of this is that the cluster was created with the wrong SSH key for the `--auth` parameter. The key supplied with the `--auth` parameter must be the same key that was used to provision the Deis servers. If you suspect this to be the issue, you'll need to `clusters:destroy` the cluster and recreate it, along with the app.

#### A Deis component fails to start
Use `fleetctl status deis-<component>.service` to get the output of the service. The most common cause of services failing to start are sporadic issues with the Docker index. The telltale sign of this is:

```console
May 12 18:24:37 deis-3 systemd[1]: Starting deis-controller...
May 12 18:24:37 deis-3 sh[6176]: 2014/05/12 18:24:37 Error: No such id: deis/controller
May 12 18:24:37 deis-3 sh[6176]: Pulling repository deis/controller
May 12 18:29:47 deis-3 sh[6176]: 2014/05/12 18:29:47 Could not find repository on any of the indexed registries.
May 12 18:29:47 deis-3 systemd[1]: deis-controller.service: control process exited, code=exited status=1
May 12 18:29:47 deis-3 systemd[1]: Failed to start deis-controller.
May 12 18:29:47 deis-3 systemd[1]: Unit deis-controller.service entered failed state.
```

We are exploring workarounds and are working with the Docker team to improve their index. In the meantime, try starting the service again with `fleetctl start deis-<component>.service`.

### Any other issues
Running into something not detailed here? Please [open an issue](https://github.com/deis/deis/issues/new) or hop into #deis and we'll help!

## License

Copyright 2014, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
