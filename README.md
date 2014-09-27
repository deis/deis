# Deis

Deis (pronounced DAY-iss) is an open source PaaS that makes it easy to deploy and manage applications on your own servers. Deis builds upon [Docker](http://docker.io/) and [CoreOS](http://coreos.com) to provide a lightweight PaaS with a [Heroku-inspired](http://heroku.com) workflow.

[![Current Release](http://img.shields.io/badge/release-v0.12.0-blue.svg)](https://github.com/deis/deis/releases/tag/v0.12.0)

![Deis Graphic](https://s3-us-west-2.amazonaws.com/deis-images/deis-graphic.png)

Deis is pre-release software. The current release is [v0.12.0](https://github.com/deis/deis/tree/v0.12.0). Until there is a stable release, we recommend you check out the latest ["master" branch](https://github.com/deis/deis) code and refer to the [latest documentation](http://docs.deis.io/en/latest/).  Upgrading from a previous Deis release? See [Upgrading Deis](http://docs.deis.io/en/latest/installing_deis/upgrading-deis/) for additional information.

# Deploying Deis

Deis is a set of Docker containers that can be deployed anywhere including public cloud, private cloud, bare metal or your workstation. Decide where you'd like to deploy Deis, then follow the deployment-specific documentation for [Rackspace](contrib/rackspace/README.md), [EC2](contrib/ec2/README.md), [DigitalOcean](contrib/digitalocean/README.md), [Google Compute Engine](contrib/gce/README.md) or [bare-metal](contrib/bare-metal/README.md) provisioning. Want to see a particular platform supported? Please open an [issue](https://github.com/deis/deis/issues/new).

Trying out Deis? Continue following these instructions for a local installation using Vagrant.

## Install prerequisites

 * Due to its nature as a distributed system, we strongly recommend using Deis with a minimum of 3 nodes even for local development and testing
 * The Deis "control plane" containers will consume approximately 2 GB of RAM across the cluster. Please be sure you have sufficient free memory before proceeding.
 * Install [Vagrant v1.6+](http://www.vagrantup.com/downloads.html) and [VirtualBox](https://www.virtualbox.org/wiki/Downloads)

Note for Ubuntu users: the VirtualBox package in Ubuntu (as of the last known release for 14.04) has some issues when running in RAM-constrained environments. Please install the latest version of VirtualBox from Oracle's website.

## Configure Discovery

Each time you spin up a new CoreOS cluster, you **must** provide a new [discovery service URL](https://coreos.com/docs/cluster-management/setup/cluster-discovery/) in the [CoreOS user-data](https://coreos.com/docs/cluster-management/setup/cloudinit-cloud-config/) file.  This URL allows hosts to find each other and perform leader election.

Automatically generate a fresh discovery URL with:

```console
$ make discovery-url
```

or manually edit [contrib/coreos/user-data](contrib/coreos/user-data) and add a unique discovery URL generated from <https://discovery.etcd.io/new>.

## Boot CoreOS

Start the CoreOS cluster on VirtualBox. From a command prompt, `cd` to the root of the Deis project code and type:

```console
$ export DEIS_NUM_INSTANCES=3
$ vagrant up
```

This instructs Vagrant to spin up 3 VMs. To be able to connect to the VMs, you must add your Vagrant-generated SSH key to the ssh-agent (`deisctl` requires the agent to have this key):

```console
$ ssh-add ~/.vagrant.d/insecure_private_key
```

## Provision Deis

Install the [deisctl utility](https://github.com/deis/deis/deisctl#installation) used to provision and operate Deis.

```console
$ curl -sSL http://deis.io/deisctl/install.sh | sudo sh
```

Export `DEISCTL_TUNNEL` so you can connect to one of the VMs using the `deisctl` client on your workstation.

```console
$ export DEISCTL_TUNNEL=172.17.8.100
```

Use `deisctl install platform` to install all Deis components across the cluster, then `deisctl start platform` to start them.

```console
$ deisctl install platform
$ deisctl start platform
```

This can take some time - the **builder** and **registry** components must download and install the beefy Heroku cedar stack.  Grab some more coffee!

Your Deis platform should be accessible at `deis.local3.deisapp.com`.  For clusters on other platforms see our guide to [Configuring DNS](http://docs.deis.io/en/latest/installing_deis/configure-dns/).

## Install the Deis Client

If you're using the latest Deis release, use `pip install --upgrade deis` to install the latest [Deis Client](https://pypi.python.org/pypi/deis/) or download [pre-compiled binaries](https://github.com/deis/deis/tree/master/client#get-started).

If you're working off master, precompiled binaries are likely out of date. You should either symlink the python file directly or build a local copy of the client:

```console
$ sudo ln -fs $(pwd)/client/deis.py /usr/local/bin/deis
```
or
```console
$ cd client && python setup.py install
```

## Register a User

Use the Deis Client to register a new user.

```console
$ deis register http://deis.local3.deisapp.com
$ deis keys:add
```

Use `deis keys:add` to add your SSH public key for `git push` access -- normally `$HOME/.ssh/id_rsa.pub`.

## Initialize a Cluster

Initialize a `dev` cluster with a list of CoreOS hosts and your CoreOS private key.

```console
$ deis clusters:create dev local3.deisapp.com --hosts=172.17.8.100 --auth=~/.vagrant.d/insecure_private_key
```

The parameters to `deis clusters:create` are:
* cluster name (`dev`) - the name used by Deis to reference the cluster
* cluster hostname (`local.3deisapp.com`) - the hostname under which apps are created, like `balancing-giraffe.local3.deisapp.com`
* cluster members (`--hosts`) - a comma-separated list of cluster members -- not necessarily all members, but at least one (for cloud providers, this is a list of the IPs like `--hosts=10.21.12.1,10.21.12.2,10.21.12.3`)
* auth SSH key (`--auth`) - the SSH private key used to provision servers -- cannot have a password (for cloud providers, this key is likely `~/.ssh/deis`)

The `dev` cluster will be used as the default cluster for future `deis` commands.

# Usage

Deis supports 3 deployment workflows:

 * Heroku Buildpacks via `git push` -- Learn more about [Using Buildpacks](http://docs.deis.io/en/latest/using_deis/using-buildpacks/)
 * Dockerfiles via `git push` -- Learn more about [Using Dockerfiles](http://docs.deis.io/en/latest/using_deis/using-dockerfiles/)
 * Docker Images via `deis pull` -- Learn more about [Using Docker Images](http://docs.deis.io/en/latest/using_deis/using-docker-images/)

As an example, we will walk through deploying a Ruby application using the Heroku Buildpack workflow.

## Prepare an Application
Clone an example Ruby application:

```console
$ git clone https://github.com/deis/example-ruby-sinatra.git
$ cd example-ruby-sinatra
```

## Create an Application
From within the application directory, create an application on Deis:

```console
$ cd example-ruby-sinatra
$ deis create
```

Use `deis create --cluster=prod` to place the app on a different cluster.  Don't like our name-generator?  Use `deis create myappname`.

## Deploy
When you created the application, a git remote for Deis was added automatically.  Deploy with `git push`.

```console
$ git push deis master
```
This will use the Deis builder to package your application as a Docker Image and automatically deploy it to the platform.
Each build creates a new release, which can be rolled back.

## Configure
Configure your application with environment variables.  Each config change also creates a new release.

```console
$ deis config:set DATABASE_URL=postgres://
```

## Test
Test the application by running your test suite inside an ephemeral Docker container.

```console
$ deis run make test
```

Use the return code to integrate with a CI system.

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

Use `deis run` to execute one-off commands and explore the deployed container.

## Testing the cluster

Integration tests and corresponding documentation can be found under the [`tests/`](tests/) folder.

## Hacking on Deis

Learn how to [hack on Deis](http://docs.deis.io/en/latest/contributing/hacking/) with a Docker-based development workflow.

## Troubleshooting

Common issues that users have run into when provisioning Deis are detailed below.

#### When running a `deisctl` command - 'Failed initializing SSH client: ssh: handshake failed: ssh: unable to authenticate'
Did you remember to add your SSH key to the ssh-agent? `ssh-add -L` should list the key you used to provision the servers. If it's not there, `ssh-add -K /path/to/your/key`.

#### When running a `deisctl` command - 'All the given peers are not reachable (Tried to connect to each peer twice and failed)'
The most common cause of this issue is that a [new discovery URL](https://discovery.etcd.io/new) wasn't generated and updated in [contrib/coreos/user-data](contrib/coreos/user-data) before the cluster was launched. Each Deis cluster must have a unique discovery URL, else there will be entries for old hosts that etcd will try and fail to connect to. Try destroying and relaunching the cluster with a fresh discovery URL.

#### Scaling an app doesn't work, and/or the app shows 'Welcome to nginx!'
This usually means the controller failed to submit jobs to the scheduler. `deisctl journal controller` will show detailed error information, but the most common cause of this is that the cluster was created with the wrong SSH key for the `--auth` parameter. The key supplied with the `--auth` parameter must be the same key that was used to provision the Deis servers. If you suspect this to be the issue, you'll need to `clusters:destroy` the cluster and recreate it, along with the app.

#### A Deis component fails to start
Use `deisctl status <component>` to view the status of the component.  You can also use `deisctl journal <component>` to tail logs for a component, or `deisctl list` to list all components.

The most common cause of services failing to start are sporadic issues with Docker Hub.  We are exploring workarounds and are working with the Docker team to improve Docker Hub reliability. In the meantime, try starting the service again with `deisctl restart <component>`.

### Any other issues
Running into something not detailed here? Please [open an issue](https://github.com/deis/deis/issues/new) or hop into [#deis](https://botbot.me/freenode/deis/) and we'll help!

## License

Copyright 2014, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
