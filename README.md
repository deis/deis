# Deis

Deis is a Django/Celery API server, Python CLI and set of [Chef cookbooks](https://github.com/opdemand/deis-cookbook) that combine to provide a Heroku-inspired application platform for public and private clouds.  Your PaaS. Your Rules.

[![Build Status](https://travis-ci.org/opdemand/deis.png)](https://travis-ci.org/opdemand/deis)
[![Coverage Status](https://coveralls.io/repos/opdemand/deis/badge.png?branch=master)](https://coveralls.io/r/opdemand/deis?branch=master)

![Deis Graphic](https://s3-us-west-2.amazonaws.com/deis-images/deis-graphic.png)

## What is Deis?

Deis is an open source PaaS that makes it easy to deploy and scale LXC containers and Chef nodes used to host applications, databases, middleware and other services. Deis leverages Chef, Docker and Heroku Buildpacks to provide a private PaaS that is lightweight and flexible.

Deis comes with out-of-the-box support for Ruby, Python, Node.js, Java, Clojure, Scala, Play, PHP, Perl, Dart and Go. However, Deis can deploy *anything* using Heroku Buildpacks, Docker images or Chef recipes.  Deis is designed to work with any cloud provider, although only EC2 is currently supported.

## Why Deis?

##### Deploy anything

Deploy a wide range of languages and frameworks with a simple `git push` using [Heroku Buildpacks](https://devcenter.heroku.com/articles/buildpacks) or [Dockerfiles](http://docs.docker.io/en/latest/use/builder/). Use custom Chef layers to deploy databases, middleware and other add-on services.

##### Control everything

Choose your hosting provider configuration. Define a [formation](http://docs.deis.io/en/latest/gettingstarted/concepts) with your own proxy and runtime layers. Retain full root access to every node. Manage your platform with a private Deis controller.

##### Scale effortlessly

Scale nodes and containers with a single command.  Node provisioning, container balancing and proxy reconfiguration are completely automated.

##### 100% Open Source

Free, transparent and easily customized. Join the open-source PaaS and DevOps community by using Deis and complimentary projects like Docker, Chef and Heroku Buildpacks.

## Getting Started

Before you get started, read about Deis core [concepts](http://docs.deis.io/en/latest/gettingstarted/concepts/) so you can answer:

 * What is a [Formation](http://docs.deis.io/en/latest/gettingstarted/concepts/#formations) and how does it relate to an application?
 * What are [Layers and Nodes](http://docs.deis.io/en/latest/gettingstarted/concepts/#layers), and how do they work with Chef?
 * How does the [Build, Release, Run](http://docs.deis.io/en/latest/gettingstarted/concepts/#build-release-run) process work?
 * How do I connect a Formation to [backing services](http://docs.deis.io/en/latest/gettingstarted/concepts/#backing-services)?

*Please note: Deis is still under active development. It should not yet be used in production.*

Follow the steps below to install your own Deis platform on EC2. To complete the installation process, you will need [Git](http://git-scm.com), [RubyGems](http://rubygems.org/pages/download), [Pip](http://www.pip-installer.org/en/latest/installing.html), the [Amazon EC2 API Tools](http://aws.amazon.com/developertools/351), [EC2 Credentials](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SettingUp_CommandLine.html#set_aws_credentials_linux) and a Chef Server with a working [Knife](http://docs.opscode.com/knife.html) client.

Don’t have a Chef Server? [Sign up for a free Hosted Chef account](https://getchef.opscode.com/signup).

### 1. Clone the Deis Repository

    $ git clone https://github.com/opdemand/deis.git
    $ cd deis

### 2. Configure the Chef Server

Deis requires a Chef Server. [Sign up for a free Hosted Chef account](https://getchef.opscode.com/signup) if you don’t have one.  You’ll also need a Ruby runtime with RubyGems in order to install the required Ruby dependencies.

	$ bundle install    # install ruby dependencies
	$ berks install     # install cookbooks into your local berkshelf
	$ berks upload      # upload cookbooks to the chef server

### 3. Provision a Deis Controller

The [Amazon EC2 API Tools](http://aws.amazon.com/developertools/351) will be used to setup basic EC2 infrastructure.  The [Knife EC2 plugin](https://github.com/opscode/knife-ec2) will be used to bootstrap the controller.

	$ contrib/provision-ec2-controller.sh

### 4. Install the Deis Client

Install the Deis client using [Pip](http://www.pip-installer.org/en/latest/installing.html).  Registration will discover SSH keys automatically and use the [standard environment variables](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SettingUp_CommandLine.html#set_aws_credentials_linux) to configure the EC2 provider.

	$ sudo pip install deis
	$ deis register http://my-deis-controller.fqdn
	username: myuser
	password: 
	email: myuser@example.com
	Registered myuser
	Logged in as myuser
	
	Found the following SSH public keys:
	1) id_rsa.pub
	Which would you like to use with Deis? 1
	Uploading /Users/myuser/.ssh/id_rsa.pub to Deis... done
	
	Found EC2 credentials: AKIAJTVXXXXXXXXXXXXX
	Import these credentials? (y/n) : y
	Uploading EC2 credentials... done

### 5. Create & Scale a Formation

Find an application you’d like to deploy, or clone [an example app](https://github.com/opdemand/example-nodejs-express).  Change into the application directory and use `deis create` to initialize a new formation in a specific EC2 region. Use the `deis layers:scale` command to provision nodes that will be dedicated to this formation.

	$ cd <my-application-repo>
	$ deis create --flavor=ec2-us-west-2
	Creating formation... done, created peachy-waxworks
	Git remote deis added
	
	Creating runtime layer... done
	Creating proxy layer... done
	
	Use deis layers:scale proxy=1 runtime=1 to scale a basic formation
	
	$ deis layers:scale proxy=1 runtime=1
	Scaling layers... but first, coffee!
	...done in 232s
	
	Use `git push deis master` to deploy to your formation

### 6. Deploy your Application

Use `git push deis master` to deploy your application.  Deis will automatically deploy Docker containers and configure Nginx proxies to route requests to your application.  To learn more, use `deis help` or browse [the documentation](http://docs.deis.io).

	$ git push deis master
	Counting objects: 146, done.
	Delta compression using up to 8 threads.
	Compressing objects: 100% (122/122), done.
	Writing objects: 100% (146/146), 21.54 KiB, done.
	Total 146 (delta 84), reused 47 (delta 22)
	       Node.js app detected
	-----> Resolving engine versions
	       Using Node.js version: 0.10.15
	       Using npm version: 1.2.30
	...
	-----> Building runtime environment
	-----> Discovering process types
	       Procfile declares types -> web
	
	-----> Compiled slug size: 4.7 MB
	       Launching... done, v2
	
	-----> peachy-waxworks deployed to Deis
	       http://ec2-54-214-143-104.us-west-2.compute.amazonaws.com ...
	
	$ curl -s http://ec2-54-214-143-104.us-west-2.compute.amazonaws.com
	Powered by Deis!

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
- Author:: Ben Grunfeld <ben@opdemand.com>

Copyright 2013, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

