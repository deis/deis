# Deis in Google Compute Engine

Let's build a Deis cluster in Google's Compute Engine!

## Prerequisites

This assumes you have a couple of items installed already:

* [deisctl](https://github.com/deis/deis/deisctl)
* `git` (Available via Homebrew or Xcode Command Line Tools)
* A clone of the Deis repository (`git clone https://github.com/deis/deis.git`)
* You are running commands from the a cloned `deis` folder (`cd deis` after cloning)

## Google

Get a few Google things squared away so we can provision VM instances.

### Google Cloud SDK

#### Install

Install the Google Cloud SDK from https://developers.google.com/compute/docs/gcutil/#install. You will then need to login with your Google Account:

```console
$ gcloud auth login
Your browser has been opened to visit:

    https://accounts.google.com/o/oauth2/auth?redirect_uri=http%3A%2F%2Flocalhost%3A8085%2F&prompt=select_account&response_type=code&client_id=22535940678.apps.googleusercontent.com&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fappengine.admin+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fbigquery+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcompute+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdevstorage.full_control+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fndev.cloudman+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcloud-platform+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fsqlservice.admin+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fprediction+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fprojecthosting&access_type=offline



You are now logged in as [youremail@gmail.com].
Your current project is [named-mason-824].  You can change this setting by running:
  $ gcloud config set project <project>
```

#### Create Project

Create a new project in the Google Developer Console (https://console.developers.google.com/project). You should get a project ID like `orbital-gantry-285` back. We'll set it as the default for the SDK tools:

```console
$ gcloud config set project orbital-gantry-285
```

#### Enable Billing

**Please note that you will begin to accrue charges once you create resources such as disks and instances**

Navigate to the project console and then the *Billing & Settings* section in the browser. Click the *Enable billing* button and fill out the form. This is needed to create resources in Google's Compute Engine.

#### Initialize Compute Engine

Google Computer Engine won't be available via the command line tools until it is initialized in the web console. Navigate to *COMPUTE* -> *COMPUTE ENGINE* -> *VM Instances* in the project console. The Compute Engine will take a moment to initialize and then be ready to create resources via `gcutil`.

### Cloud Init

Create your cloud init file using the Deis `contrib/gce/create-gce-user-data` script and a new etcd discovery URL. First install the PyYAML:

```console
$ sudo pip install pyyaml
Downloading/unpacking pyyaml
  Downloading PyYAML-3.11.tar.gz (248kB): 248kB downloaded
  Running setup.py (path:/private/tmp/pip_build_root/pyyaml/setup.py) egg_info for package pyyaml

Installing collected packages: pyyaml
  Running setup.py install for pyyaml

  ...

Successfully installed pyyaml
Cleaning up...
```

Then navigate to the `contrib/gce` directory:

```console
$ cd contrib/gce
```

Finally, create the `gce-user-data` file:

```console
$ ./create-gce-user-data $(curl -s https://discovery.etcd.io/new)
Wrote file /Users/andy/Projects/deis/contrib/gce/gce-user-data with discovery URL https://discovery.etcd.io/ca4af6f7534f52055e889fddc8c9c4ae
```

We should have a `gce-user-data` file ready to launch CoreOS nodes with.

### Launch Instances

Create a SSH key that we will use for Deis host communication:

```console
$ ssh-keygen -q -t rsa -f ~/.ssh/deis -N '' -C deis
```

Create some persistent disks to use for `/var/lib/docker`. The default root partition of CoreOS is only around 4 GB and not enough for storing Docker images and instances. The following creates 3 disks sized at 32 GB:

```console
$ gcutil adddisk --zone us-central1-a --size_gb 32 cored1 cored2 cored3

Table of resources:

+--------+---------------+--------+---------+
| name   | zone          | status | size-gb |
+--------+---------------+--------+---------+
| cored1 | us-central1-a | READY  |      32 |
+--------+---------------+--------+---------+
| cored2 | us-central1-a | READY  |      32 |
+--------+---------------+--------+---------+
| cored3 | us-central1-a | READY  |      32 |
+--------+---------------+--------+---------+
```

Launch 3 instances using `coreos-alpha-444-0-0-v20140919` image. You can choose another starting CoreOS image from the listing output of `gcloud compute images list`:

```console
$ for num in 1 2 3; do gcutil addinstance --image projects/coreos-cloud/global/images/coreos-alpha-444-0-0-v20140919 --persistent_boot_disk --zone us-central1-a --machine_type n1-standard-2 --tags deis --metadata_from_file user-data:gce-user-data --disk cored${num},deviceName=coredocker --authorized_ssh_keys=core:~/.ssh/deis.pub,core:~/.ssh/google_compute_engine.pub core${num}; done

Table of resources:

+-------+---------------+--------------+---------------+---------+
| name  | network-ip    | external-ip  | zone          | status  |
+-------+---------------+--------------+---------------+---------+
| core1 | 10.240.33.107 | 23.236.59.66 | us-central1-a | RUNNING |
+-------+---------------+--------------+---------------+---------+
| core2 | 10.240.94.33  | 108.59.80.17 | us-central1-a | RUNNING |
+-------+---------------+--------------+---------------+---------+
| core3 | 10.240.28.163 | 108.59.85.85 | us-central1-a | RUNNING |
+-------+---------------+--------------+---------------+---------+
```

### Load Balancing

We will need to load balance the Deis routers so we can get to Deis services (controller and builder) and our applications.

```console
$ gcutil addhttphealthcheck basic-check --request_path /health-check
$ gcutil addtargetpool deis --health_checks basic-check --region us-central1 --instances core1,core2,core3
$ gcutil addforwardingrule deisapp --region us-central1 --target_pool deis

Table of resources:

+---------+-------------+--------------+
| name    | region      | ip           |
+---------+-------------+--------------+
| deisapp | us-central1 | 23.251.153.6 |
+---------+-------------+--------------+
```

Note the forwarding rule external IP address. We will use it as the Deis login endpoint in a future step. Now allow the ports on the CoreOS nodes:

```console
$ gcutil addfirewall deis-router --target_tags deis --allowed "tcp:80,tcp:2222"
```

### DNS

We can create DNS records in Google Cloud DNS using the `gcloud` utility. In our example we will be using the domain name `deisdemo.io`. Create the zone:

```console
$ gcloud dns managed-zone create --dns_name deisdemo.io. --description "Example Deis cluster domain name" deisdemoio
Creating {'dnsName': 'deisdemo.io.', 'name': 'deisdemoio', 'description':
'Example Deis cluster domain name'} in eco-theater-654

Do you want to continue (Y/n)?  Y

{
    "creationTime": "2014-07-28T00:01:45.835Z",
    "description": "Example Deis cluster domain name",
    "dnsName": "deisdemo.io.",
    "id": "1374035518570040348",
    "kind": "dns#managedZone",
    "name": "deisdemoio",
    "nameServers": [
        "ns-cloud-d1.googledomains.com.",
        "ns-cloud-d2.googledomains.com.",
        "ns-cloud-d3.googledomains.com.",
        "ns-cloud-d4.googledomains.com."
    ]
}
```

Note the `nameServers` array from the JSON output. We will need to setup our upstream domain name servers to these.

Now edit the zone to add the Deis endpoint and wildcard DNS:

```console
$ gcloud dns records --zone deisdemoio edit
{
    "additions": [
        {
            "kind": "dns#resourceRecordSet",
            "name": "deisdemo.io.",
            "rrdatas": [
                "ns-cloud-d1.googledomains.com. dns-admin.google.com. 2 21600 3600 1209600 300"
            ],
            "ttl": 21600,
            "type": "SOA"
        }
    ],
    "deletions": [
        {
            "kind": "dns#resourceRecordSet",
            "name": "deisdemo.io.",
            "rrdatas": [
                "ns-cloud-d1.googledomains.com. dns-admin.google.com. 1 21600 3600 1209600 300"
            ],
            "ttl": 21600,
            "type": "SOA"
        }
    ]
}
```

You will want to add two records as JSON objects. Here is an example edit for the two A record additions:

```json
{
    "additions": [
        {
            "kind": "dns#resourceRecordSet",
            "name": "deisdemo.io.",
            "rrdatas": [
                "ns-cloud-d1.googledomains.com. dns-admin.google.com. 2 21600 3600 1209600 300"
            ],
            "ttl": 21600,
            "type": "SOA"
        },
        {
            "kind": "dns#resourceRecordSet",
            "name": "deis.deisdemo.io.",
            "rrdatas": [
                "23.251.153.6"
            ],
            "ttl": 21600,
            "type": "A"
        },
        {
            "kind": "dns#resourceRecordSet",
            "name": "*.dev.deisdemo.io.",
            "rrdatas": [
                "23.251.153.6"
            ],
            "ttl": 21600,
            "type": "A"
        }
    ],
    "deletions": [
        {
            "kind": "dns#resourceRecordSet",
            "name": "deisdemo.io.",
            "rrdatas": [
                "ns-cloud-d1.googledomains.com. dns-admin.google.com. 1 21600 3600 1209600 300"
            ],
            "ttl": 21600,
            "type": "SOA"
        }
    ]
}
```

## Deis

Time to install Deis!

### Install

We cloned the Deis repository in the prerequisites. In this example we will be deploying version `0.12.0`:

```console
$ git checkout v0.12.0
Note: checking out 'v0.12.0'.

You are in 'detached HEAD' state. You can look around, make experimental
changes and commit them, and you can discard any commits you make in this
state without impacting any branches by performing another checkout.

If you want to create a new branch to retain commits you create, you may
do so (now or later) by using -b with the checkout command again. Example:

  git checkout -b new_branch_name

HEAD is now at 64708ab... chore(docs): update CLI versions and download links
```

Then install the CLI:

```console
$ sudo pip install --upgrade ./client
```

### Setup

The `DEISCTL_TUNNEL` environment variable provides a SSH gateway to use in with `deisctl`. Use the public IP address for one of the CoreOS nodes we deployed earlier:

```shell
export DEISCTL_TUNNEL=23.236.59.66
```

Verify the CoreOS cluster is operational and that we can communicate with the cluster:

```console
$ deisctl list
MACHINE		IP		METADATA
```

Now we can bootstrap the Deis containers:

```shell
deisctl install platform && deisctl start platform
```

This operation will take a while as all the Deis systemd units are loaded into the CoreOS cluster and the Docker images are pulled down. Grab some iced tea!

Verify that all the units are active after the operation completes:

```console
$ deisctl list
UNIT                        MACHINE                     LOAD    ACTIVE  SUB
deis-builder@1.service      dea53588.../172.17.8.100    loaded  active  running
deis-cache@1.service        dea53588.../172.17.8.100    loaded  active  running
deis-controller@1.service   dea53588.../172.17.8.100    loaded  active  running
deis-database-data.service  dea53588.../172.17.8.100    loaded  active  exited
deis-database@1.service     dea53588.../172.17.8.100    loaded  active  running
deis-logger-data.service    dea53588.../172.17.8.100    loaded  active  exited
deis-logger@1.service       dea53588.../172.17.8.100    loaded  active  running
deis-registry-data.service  dea53588.../172.17.8.100    loaded  active  exited
deis-registry@1.service     dea53588.../172.17.8.100    loaded  active  running
deis-router@1.service       dea53588.../172.17.8.100    loaded  active  running
```

Everything looks good! Register the admin user. The first user added to the system becomes the admin:

```console
$ deis register http://deis.deisdemo.io
username: andyshinn
password:
password (confirm):
email: andys@andyshinn.as
Registered andyshinn
Logged in as andyshinn
```

You are now registered and logged in. Create a new cluster named `dev` to run applications under. You could name this cluster something other than `dev`. We only use it as an example to illustrate a cluster can be restricted to certain CoreOS nodes. The hosts are the internal IP addresses of the CoreOS nodes:

```console
$ deis clusters:create dev deisdemo.io --hosts 10.240.33.107,10.240.94.33,10.240.28.163 --auth ~/.ssh/deis
Creating cluster... done, created dev
```

Add your SSH key so you can publish applications:

```console
$ deis keys:add
Found the following SSH public keys:
1) id_rsa.pub andy
Which would you like to use with Deis? 1
Uploading andy to Deis...done
```

### Applications

Creating an application requires that application be housed under git already. Clone an example application and deploy it:

```console
$ git clone https://github.com/deis/example-ruby-sinatra.git
Cloning into 'example-ruby-sinatra'...
remote: Counting objects: 98, done.
remote: Compressing objects: 100% (50/50), done.
remote: Total 98 (delta 42), reused 97 (delta 42)
Unpacking objects: 100% (98/98), done.
Checking connectivity... done.
$ cd example-ruby-sinatra
$ deis
$ deis create
Creating application... done, created breezy-frosting
Git remote deis added
```

Time to push:

```console
$ git push deis master
Counting objects: 98, done.
Delta compression using up to 8 threads.
Compressing objects: 100% (92/92), done.
Writing objects: 100% (98/98), 20.95 KiB | 0 bytes/s, done.
Total 98 (delta 42), reused 0 (delta 0)
-----> Ruby app detected
-----> Compiling Ruby/Rack
-----> Using Ruby version: ruby-1.9.3
-----> Installing dependencies using 1.6.3
       Running: bundle install --without development:test --path vendor/bundle --binstubs vendor/bundle/bin -j4 --deployment
       Don't run Bundler as root. Bundler can ask for sudo if it is needed, and
       installing your bundle as root will break this application for all non-root
       users on this machine.
       Fetching gem metadata from http://rubygems.org/..........
       Fetching additional metadata from http://rubygems.org/..
       Using bundler 1.6.3
       Installing rack 1.5.2
       Installing tilt 1.3.6
       Installing rack-protection 1.5.0
       Installing sinatra 1.4.2
       Your bundle is complete!
       Gems in the groups development and test were not installed.
       It was installed into ./vendor/bundle
       Bundle completed (5.72s)
       Cleaning up the bundler cache.
-----> Discovering process types
       Procfile declares types -> web
       Default process types for Ruby -> rake, console, web
-----> Compiled slug size is 12M
remote: -----> Building Docker image
remote: Sending build context to Docker daemon 11.77 MB
remote: Sending build context to Docker daemon
remote: Step 0 : FROM deis/slugrunner
remote:  ---> f607bc8783a5
remote: Step 1 : RUN mkdir -p /app
remote:  ---> Running in dd1cb10534c0
remote:  ---> 3151b07f7623
remote: Removing intermediate container dd1cb10534c0
remote: Step 2 : ADD slug.tgz /app
remote:  ---> b86143c577ae
remote: Removing intermediate container 63dca22b29d6
remote: Step 3 : ENTRYPOINT ["/runner/init"]
remote:  ---> Running in 43c572eacc69
remote:  ---> 6eeace9fea7e
remote: Removing intermediate container 43c572eacc69
remote: Successfully built 6eeace9fea7e
remote: -----> Pushing image to private registry
remote:
remote:        Launching... done, v2
remote:
remote: -----> breezy-frosting deployed to Deis
remote:        http://breezy-frosting.dev.deisdemo.io
remote:
remote:        To learn more, use `deis help` or visit http://deis.io
remote:
To ssh://git@deis.deisdemo.io:2222/breezy-frosting.git
 * [new branch]      master -> master
```

Your application will now be built and run inside the Deis cluster! After the application is pushed it should be running at http://breezy-frosting.dev.deisdemo.io. Check the application information:

```console
$ deis apps:info
=== breezy-frosting Application
{
  "updated": "2014-07-28T00:35:45.528Z",
  "uuid": "fd926c94-5b65-48e8-8afe-7ac547c12bd6",
  "created": "2014-07-28T00:33:35.346Z",
  "cluster": "dev",
  "owner": "andyshinn",
  "id": "breezy-frosting",
  "structure": "{\"web\": 1}"
}

=== breezy-frosting Processes

--- web:
web.1 up (v2)

=== breezy-frosting Domains
No domains
```

Can we connect to the application?

```console
$ curl -s http://breezy-frosting.dev.deisdemo.io
Powered by Deis!
```

It works! Enjoy your Deis cluster in Google Compute Engine!
