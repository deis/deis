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
 * How do I connect an application to [backing services](http://docs.deis.io/en/latest/gettingstarted/concepts/#backing-services)?

Follow the steps below to install your own Deis platform on EC2. To complete the installation process, you will need [Git](http://git-scm.com), [RubyGems](http://rubygems.org/pages/download), [Pip](http://www.pip-installer.org/en/latest/installing.html), the [Amazon EC2 API Tools](http://aws.amazon.com/developertools/351), [EC2 Credentials](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SettingUp_CommandLine.html#set_aws_credentials_linux) and a Chef Server with a working [Knife](http://docs.opscode.com/knife.html) client.

*Please note: Deis is still under active development. It should not yet be used in production.*

### 1. Clone the Deis Repository

```bash
$ git clone https://github.com/opdemand/deis.git
$ cd deis
```

Cloning the default master branch will provide you with the latest development version of Deis.  If you want to deploy the latest stable release, make sure you checkout the most recent tag using ``git checkout vX.Y.Z``.

### 2. Configure the Chef Server

Deis requires a Chef Server. [Sign up for a free Hosted Chef account](https://getchef.opscode.com/signup) if you don’t have one.  You’ll also need a Ruby runtime with RubyGems in order to install the required Ruby dependencies.

```bash
$ bundle install    # install ruby dependencies
$ berks install     # install cookbooks into your local berkshelf
$ berks upload      # upload cookbooks to the chef server
```

### 3. Provision a Deis Controller

The [Amazon EC2 API Tools](http://aws.amazon.com/developertools/351) will be used to setup basic EC2 infrastructure.  The [Knife EC2 plugin](https://github.com/opscode/knife-ec2) will be used to bootstrap the controller.

	$ contrib/ec2/provision-ec2-controller.sh

Once the `deis-controller` node exists on the Chef server, you *must* log in to the WebUI add deis-controller to the `admins` group.  This is required so the controller can delete node and client records during future scaling operations.

### 4. Install the Deis Client

Install the Deis client using [Pip](http://www.pip-installer.org/en/latest/installing.html) (for latest stable) or by linking `<repo>/client/deis.py` to `/usr/local/bin/deis` (for dev version).  Registration will discover SSH keys automatically and use the [standard environment variables](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/SettingUp_CommandLine.html#set_aws_credentials_linux) to configure the EC2 provider.

```bash
$ sudo pip install deis
$ deis register http://my-deis-controller.fqdn
username: myuser
password:
password (confirm):
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
```

### 5. Create & Scale a Formation

Use the Deis client to create a new formation named "dev" that
has a default layer that serves as both runtime (hosts containers)
and proxy (routes traffic to containers).  Scale the default layer
up to one node.

```bash
$ deis formations:create dev --flavor=ec2-us-west-2 --domain=deisapp.com
Creating formation... done, created dev
Creating runtime layer... done in 1s

Use `deis nodes:scale dev runtime=1` to scale a basic formation

$ deis nodes:scale dev runtime=1
Scaling nodes... but first, coffee!
...done in 251s

Use `deis create --formation=dev` to create an application
```

### 6. Deploy & Scale an Application

Change into your application directory and use  ``deis create --formation=dev``
to create a new application attached to the dev formation.

To deploy the application, use `git push deis master`.  Deis will automatically deploy Docker containers and configure Nginx proxies to route requests to your application.

Once your application is deployed, use ``deis scale web=4`` to
scale up web containers.  You can also use ``deis logs`` to view
aggregated application logs, or ``deis run`` to run admin
commands inside your application.

To learn more, use `deis help` or browse [the documentation](http://docs.deis.io).

```bash
$ deis create --formation=dev
Creating application... done, created peachy-waxworks
Git remote deis added

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
       http://peachy-waxworks.deisapp.com ...

$ curl -s http://peachy-waxworks.deisapp.com
Powered by Deis!

$ deis scale web=4
Scaling containers... but first, coffee!
done in 12s

=== peachy-waxworks Containers

--- web: `node server.js`
web.1 up 2013-09-23T19:02:30.745Z (dev-runtime-1)
web.2 up 2013-09-23T19:36:48.741Z (dev-runtime-1)
web.3 up 2013-09-23T19:36:48.758Z (dev-runtime-1)
web.4 up 2013-09-23T19:36:48.771Z (dev-runtime-1)
```

## Credits

Deis stands on the shoulders of leading open source technologies:

  * [Chef](http://www.opscode.com/)
  * [Docker](http://www.docker.io/)
  * [Django](https://www.djangoproject.com/)
  * [Celery](http://www.celeryproject.org/)
  * [Heroku Buildpacks](https://devcenter.heroku.com/articles/buildpacks)
  * [Buildstep](https://github.com/progrium/buildstep)
  * [Gitosis](https://github.com/opdemand/gitosis)

## License and Authors

- Author:: Gabriel Monroy <gabriel@opdemand.com>
- Author:: Matt Boersma <matt@opdemand.com>
- Author:: Ben Grunfeld <ben@opdemand.com>

Copyright 2013, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

