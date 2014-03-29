# Deis

Deis is an open source PaaS that makes it easy to deploy, scale and manage Docker containers used to host applications and services. Deis builds upon Docker and CoreOS to provide a private PaaS that is lightweight and flexible.

[![Build Status](https://travis-ci.org/opdemand/deis.png?branch=master)](https://travis-ci.org/opdemand/deis)
[![Coverage Status](https://coveralls.io/repos/opdemand/deis/badge.png?branch=master)](https://coveralls.io/r/opdemand/deis?branch=master)

![Deis Graphic](https://s3-us-west-2.amazonaws.com/deis-images/deis-graphic.png)

# Installation

Deis is a set of Docker containers that can be deployed anywhere including public cloud, private cloud, bare metal or your workstation.  You will need Docker and a CoreOS cluster to get started.

## Run Deis

Build Deis and run the `deis/deis` Docker image.

```
make build
make run
```

## Install the Deis Client
Use `pip` to install the latest Deis Client, or download pre-compiled binares.

```
pip install deis
```

## Register a User
Use the Deis Client to register a new user.

```
deis register http://localhost:8000
```

## Initalize a Cluster

Initalize a `dev` cluster with a list of CoreOS hosts and your CoreOS private key.

```
deis clusters:create dev deisapp.com --hosts=coreos-host1,coreos-host2,coreos-host3 --auth=~/.ssh/coreos
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

#### From a Docker Registry

You can also push builds directly from the Docker Index or a private Docker registry.  First, build and push the Docker images as you normally do:

```
docker build -t gabrtv/example
docker push gabrtv/example
```

Then use `deis push` to deploy.

```
deis push gabrtv/example
```

Use the fully qualfied image path to push from a private registry: `deis push registry.local:5000/gabrtv/example`

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
deis scale 8
```

To scale by process type, use `deis scale web=8 worker=2` .  Just make sure the process types can be run via `/start web`.

## Publish
Publish your application via each cluster's integrated router.

```
deis publish 8080/http
deis domain myapp.example.com
```

Use `deis open` to pop into a browser pointed at your application.  Use `--sslCert=cert.pem --sslKey=key.pem` to secure the connection with SSL/TLS.

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
