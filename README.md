# Deis

Deis is a Django/Celery API server, Python CLI and set of [Chef cookbooks](https://github.com/opdemand/deis-cookbook) that combine to provide a Heroku-inspired application platform for public and private clouds.  Your PaaS. Your Rules.

[![Build Status](https://travis-ci.org/opdemand/deis.png)](https://travis-ci.org/opdemand/deis)
[![Coverage Status](https://coveralls.io/repos/opdemand/deis/badge.png?branch=master)](https://coveralls.io/r/opdemand/deis?branch=master)

## What is Deis?
Deis is an open source PaaS that makes it easy to deploy and scale LXC containers and Chef nodes used to host applications, databases, middleware and other services. Deis leverages Chef, Docker and Heroku Buildpacks to provide a private PaaS that is lightweight and flexible.Deis comes with out-of-the-box support for Ruby, Python, Node.js, Java, Clojure, Scala, Play, PHP, Perl, Dart and Go.  However, Deis can deploy *anything* using Heroku Buildpacks, Docker images or Chef recipes.  Deis is designed to work with any cloud provider, although only EC2 is currently supported.## Why Deis?##### Deploy anything

Deploy a wide range of languages and frameworks with a simple "git push" using [Heroku Buildpacks](https://devcenter.heroku.com/articles/buildpacks) or [Dockerfiles](http://docs.docker.io/en/latest/use/builder/).  Use custom Chef layers to deploy databases, middleware and other add-on services.
##### Control everything
Choose your hosting providers.  Define a "formation" with custom proxy and runtime layers.  Scale nodes and containers independently.  Manage the entire platform with a private Deis controller.
##### Scale effortlessly
Scale nodes and containers with a single command.  Provisioning is transparent, container formations are rebalanced automatically and proxies are updated to re-route traffic without downtime.
##### 100% Open Source
Free, transparent and easily customized.  Join the open-source PaaS and DevOps community by using Deis and complimentary projects like Docker, Chef and Heroku Buildpacks.## Getting Started
Coming Soon!

## Credits

Deis rests on the shoulders of leading open source technologies:

  * [Chef](http://www.opscode.com/)
  * [Docker](http://www.docker.io/)
  * [Django](https://www.djangoproject.com/)
  * [Celery](http://www.celeryproject.org/)
  * [Heroku](https://devcenter.heroku.com/articles/buildpacks)
  * [Buildstep](https://github.com/progrium/buildstep)
  * [Gitosis](https://github.com/opdemand/gitosis)

## License and Authors

- Author:: Gabriel Monroy <gabriel@opdemand.com>
- Author:: Matt Boersma <matt@opdemand.com>

Copyright 2013, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
