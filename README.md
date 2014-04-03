# Deis

Deis is an open source PaaS that makes it easy to deploy, scale and manage Docker containers used to host applications and services. Deis builds upon Docker and CoreOS to provide a private PaaS that is lightweight and flexible.

[![Build Status](https://travis-ci.org/opdemand/deis.png?branch=master)](https://travis-ci.org/opdemand/deis)
[![Coverage Status](https://coveralls.io/repos/opdemand/deis/badge.png?branch=master)](https://coveralls.io/r/opdemand/deis?branch=master)

![Deis Graphic](https://s3-us-west-2.amazonaws.com/deis-images/deis-graphic.png)

# Installation

Deis is a set of Docker containers that can be deployed anywhere including public cloud, private cloud, bare metal or your workstation.  You will need Docker and Vagrant to get started.

## Boot CoreOS

Start a CoreOS virtual machine on VirtualBox.

```
vagrant up
```

Export some environment variables so you can connect to the VM using the `docker` and `fleetctl` clients on your workstation.

```
export DOCKER_HOST=tcp://172.17.8.100:4243
export FLEETCTL_TUNNEL=172.17.8.100
```

## Build Deis

Use `make build` to assemble all of the Deis components from Dockerfiles.  Grab some coffee while it builds the images on the CoreOS VM (it can take a while).

```
make build
```

## Run Deis

Use `make run` to start all Deis containers and attach to their log output.

```
make run
```

## Install the Deis Client
Use `pip` to install the latest Deis Client, download pre-compiled binares, or symlink `client/deis.py` to use the latest version.

```
ln -fs $(pwd)/client/deis.py /usr/local/bin/deis
```

## Register a User
Use the Deis Client to register a new user.

```
deis register http://local.deisapp.com:8000
deis keys:add
```

Use `deis keys:add` to add your SSH public key for `git push` access.

## Initalize a Cluster

Initalize a `dev` cluster with a list of CoreOS hosts and your CoreOS private key.

```
deis clusters:create dev local.deisapp.com --hosts=local.deisapp.com --auth=~/.vagrant.d/insecure_private_key
```

The `dev` cluster will be used as the default cluster for future `deis` commands.

# Usage

## Create an Application
Create an application on the default `dev` cluster.

```
deis create
```

Use `deis create --cluster=prod` to place the app on a different cluster.  Don't like our name-generator?  Use `deis create myappname`.

## Push
Push builds of your application from your local git repository or from a Docker Registry.  Each build creates a new release, which can be rolled back.

#### From a Git Repository
When you created the application, a git remote for Deis was added automatically.

```
git push deis master
```
This will use the Deis builder to package your application as a Docker Image and deploy it on your application's cluster.

## Configure
Configure your application with environment variables.  Each config change also creates a new release.

```
deis config:set DATABASE_URL=postgres://
```

Coming soon: Use the integrated ETCD namespace for service discovery between applications on the same cluster.

## Test
Test your application by running commands inside an ephemeral Docker container.

```
deis run make test
```

To integrate with your CI system, check the return code.

## Scale
Scale containers horizontally with ease.

```
deis scale web=8
```

## Debug
Access to aggregated logs makes it easy to troubleshoot problems with your application.

```
deis logs
```

Use `deis run` to execute one-off commands and explore the deployed container.  Coming soon: `deis attach` to jump into a live container.

## License

Copyright 2014, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
