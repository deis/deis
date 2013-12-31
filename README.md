# Deis

Deis is a Django/Celery API server, Python CLI and set of [Chef cookbooks](https://github.com/opdemand/deis-cookbook) that combine to provide a Heroku-inspired application platform for public and private clouds.  Your PaaS. Your Rules.

[![Build Status](https://travis-ci.org/opdemand/deis.png)](https://travis-ci.org/opdemand/deis)
[![Coverage Status](https://coveralls.io/repos/opdemand/deis/badge.png?branch=master)](https://coveralls.io/r/opdemand/deis?branch=master)

![Deis Graphic](https://s3-us-west-2.amazonaws.com/deis-images/deis-graphic.png)

## What is Deis?

Deis is an open source PaaS that makes it easy to deploy and scale LXC containers and Chef nodes used to host applications, databases, middleware and other services. Deis leverages Chef, Docker and Heroku Buildpacks to provide a private PaaS that is lightweight and flexible.

Deis comes with out-of-the-box support for Ruby, Python, Node.js, Java, Clojure, Scala, Play, PHP, Perl, Dart and Go. However, Deis can deploy *anything* using Heroku Buildpacks, Docker images or Chef recipes.  Deis can be deployed on any system including every public cloud, private cloud or bare metal.

## Why Deis?

##### Deploy anything

Deploy a wide range of languages and frameworks with a simple `git push` using [Heroku Buildpacks](https://devcenter.heroku.com/articles/buildpacks) or (coming soon) [Dockerfiles](http://docs.docker.io/en/latest/use/builder/). Use custom Chef layers to deploy databases, middleware and other add-on services.

##### Control everything

Choose your hosting provider configuration. Define a [formation](http://docs.deis.io/en/latest/gettingstarted/concepts) with your own proxy and runtime layers. Retain full root access to every node. Manage your platform with a private Deis controller.

##### Scale effortlessly

Scale nodes and containers with a single command.  Node provisioning, container balancing and proxy reconfiguration are completely automated.

##### 100% Open Source

Free, transparent and easily customized. Join the open-source PaaS and DevOps community by using Deis and complimentary projects like Docker, Chef and Heroku Buildpacks.

## Getting Started

First read about Deis core [concepts](http://docs.deis.io/en/latest/gettingstarted/concepts/) so you can answer:

 * What is a [Formation](http://docs.deis.io/en/latest/gettingstarted/concepts/#formations) and how does it relate to an application?
 * What are [Layers and Nodes](http://docs.deis.io/en/latest/gettingstarted/concepts/#layers), and how do they work with Chef?
 * How does the [Build, Release, Run](http://docs.deis.io/en/latest/gettingstarted/concepts/#build-release-run) process work?
 * How do I connect an application to [backing services](http://docs.deis.io/en/latest/gettingstarted/concepts/#backing-services)?

Next, choose which cloud provider should host your Deis controller.
Regardless of whether your controller is hosted with
[Amazon EC2](http://docs.deis.io/en/latest/installation/ec2/),
[Rackspace](http://docs.deis.io/en/latest/installation/rackspace/), or
[DigitalOcean](http://docs.deis.io/en/latest/installation/digitalocean/),
it can create and manage nodes on any of those cloud providers.

Proceed to the [Operations Guide](http://docs.deis.io/en/latest/operations/)
documentation to start building your own private PaaS.

## Credits

Deis stands on the shoulders of leading open source technologies:

  * [Chef](http://www.opscode.com/)
  * [Docker](http://www.docker.io/)
  * [Django](https://www.djangoproject.com/)
  * [Celery](http://www.celeryproject.org/)
  * [Heroku Buildpacks](https://devcenter.heroku.com/articles/buildpacks)
  * [Slugbuilder](https://github.com/flynn/slugbuilder) and [slugrunner](https://github.com/flynn/slugrunner)
  * [Gitosis](https://github.com/opdemand/gitosis)

## License and Authors

- Author:: Gabriel Monroy <gabriel@opdemand.com>
- Author:: Matt Boersma <matt@opdemand.com>
- Author:: Ben Grunfeld <ben@opdemand.com>

Copyright 2013, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

