# Deis Control Utility

`deisctl` is a command-line utility used to provision and operate a Deis cluster.

## Installation

To install `deisctl` on Linux or Mac OS X, run this command:

```console
$ curl -sSL http://deis.io/deisctl/install.sh | sudo sh
```

The installer puts `deisctl` in */usr/local/bin* and downloads current Deis unit files
to */var/lib/deis/units* one time.

To change installation options, save the installer directly from one of these links:

[![Download for Linux](http://img.shields.io/badge/download-Linux-brightgreen.svg?style=flat)](https://s3-us-west-2.amazonaws.com/opdemand/deisctl-0.13.0-dev-linux-amd64.run)
[![Download for Mac OS X](http://img.shields.io/badge/download-Mac%20OS%20X-brightgreen.svg?style=flat)](https://s3-us-west-2.amazonaws.com/opdemand/deisctl-0.13.0-dev-darwin-amd64.run)

Then run the downloaded file as a shell script. Append `--help` to see what options
are available.

If you want to install from source, clone the repository and run

```console
$ godep get .
```

Then, export the `DEISCTL_UNITS` environment variable so deisctl can find the units:

```console
$ export DEISCTL_UNITS="$PATH_TO_DEISCTL/units"
```

## Remote Configuration

While `deisctl` can be used locally on a CoreOS host, it is extremely useful as a tool
for remote administration.  This requires an SSH tunnel to one of your CoreOS hosts.

Test password-less SSH connectivity to a CoreOS host:

```console
$ ssh core@172.17.8.100 hostname
deis-1
```

Export the `DEISCTL_TUNNEL` environment variable:

```console
$ export DEISCTL_TUNNEL=172.17.8.100
```

## Provision a Deis Platform

The `deisctl install platform` command will schedule all of the Deis platform
units. `deisctl start platform` activates these units.

```console
$ deisctl install platform
● ▴ ■
■ ● ▴ Installing Deis...
▴ ■ ●

Scheduling data containers...
deis-database-data.service: loaded
deis-registry-data.service: loaded
deis-logger-data.service: loaded
Data containers scheduled.
Scheduling service containers...
deis-database@1.service: loaded
deis-cache@1.service: loaded
deis-logger@1.service: loaded
deis-registry@1.service: loaded
deis-controller@1.service: loaded
deis-builder@1.service: loaded
deis-router@1.service: loaded
Service containers scheduled.
Deis installed.
Please run `deisctl start platform` to boot up Deis.

$ deisctl start platform
● ▴ ■
■ ● ▴ Starting Deis...
▴ ■ ●

Launching data containers...
deis-database-data.service: exited
deis-registry-data.service: exited
deis-logger-data.service: exited
Data containers launched.
Launching service containers...
deis-logger@1.service: running
deis-cache@1.service: running
deis-router@1.service: running
deis-database@1.service: running
deis-controller@1.service: running
deis-registry@1.service: running
deis-builder@1.service: running
Deis started.
```

Note that the default start command activates 1 of each component.
You can scale components with `deisctl scale router=3`, for example.
The router is the only component that _currently_ scales beyond 1 unit.

You can also use the `deisctl uninstall` command to destroy platform units:

```console
$ deisctl uninstall platform

Destroying service containers...
deis-database@1.service: inactive
deis-cache@1.service: inactive
deis-logger@1.service: inactive
deis-registry@1.service: inactive
deis-controller@1.service: inactive
deis-builder@1.service: inactive
deis-router@1.service: inactive
Done.
```

To uninstall a specific component, use `deisctl uninstall router`.

Note that uninstalling platform units will _not_ remove the data units or underlying
data containers.  Data must be destroyed manually.

## Usage

The `deisctl` tool provides a number of other commands, including:

 * `deisctl list` - list Deis platform components
 * `deisctl status <component>` - retrieve Systemd status of a component
 * `deisctl journal <component>` - retrieve Systemd journal output
 * `deisctl start <component>` - start a platform component
 * `deisctl stop <component>` - stop a platform component
 * `deisctl install <component>` - install a single platform component
 * `deisctl uninstall <component>` - uninstall a single platform component
 * `deisctl scale <component>=<num>` - scale a component to the target number of units
 * `deisctl refresh-units` - download latest unit files

## Usage Examples

```console
$ deisctl list
UNIT				MACHINE				LOAD	ACTIVE	SUB
deis-builder@1.service		f936b7a5.../172.17.8.100	loaded	active	running
deis-cache@1.service		f936b7a5.../172.17.8.100	loaded	active	running
deis-controller@1.service	f936b7a5.../172.17.8.100	loaded	active	running
deis-database-data.service	f936b7a5.../172.17.8.100	loaded	active	exited
deis-database@1.service		f936b7a5.../172.17.8.100	loaded	active	running
deis-logger-data.service	f936b7a5.../172.17.8.100	loaded	active	exited
deis-logger@1.service		f936b7a5.../172.17.8.100	loaded	active	running
deis-registry-data.service	f936b7a5.../172.17.8.100	loaded	active	exited
deis-registry@1.service		f936b7a5.../172.17.8.100	loaded	active	running
deis-router@1.service		f936b7a5.../172.17.8.100	loaded	active	running
```

```console
$ deisctl status controller
● deis-controller@1.service - deis-controller
   Loaded: loaded (/run/fleet/units/deis-controller@1.service; linked-runtime)
   Active: active (running) since Mon 2014-08-25 22:56:50 UTC; 15min ago
  Process: 22969 ExecStartPre=/bin/sh -c docker inspect deis-controller >/dev/null && docker rm -f deis-controller || true (code=exited, status=0/SUCCESS)
  Process: 22945 ExecStartPre=/bin/sh -c IMAGE=`/run/deis/bin/get_image /deis/controller`; docker history $IMAGE >/dev/null || docker pull $IMAGE (code=exited, status=0/SUCCESS)
 Main PID: 22979 (sh)
   CGroup: /system.slice/system-deis\x2dcontroller.slice/deis-controller@1.service
           ├─22979 /bin/sh -c IMAGE=`/run/deis/bin/get_image /deis/controller` && docker run --name deis-controller --rm -p 8000:8000 -e PUBLISH=8000 -e HOST=$COREOS_PRIVATE_IPV4 --volumes-from=deis-logger $IMAGE
           └─22999 docker run --name deis-controller --rm -p 8000:8000 -e PUBLISH=8000 -e HOST=172.17.8.100 --volumes-from=deis-logger deis/controller:latest

Aug 25 22:57:07 deis-1 sh[22979]: [2014-08-25 16:57:07,959: INFO/MainProcess] Connected to redis://172.17.8.100:6379/0
Aug 25 22:57:07 deis-1 sh[22979]: 2014-08-25 16:57:07 [121] [INFO] Booting worker with pid: 121
Aug 25 22:57:07 deis-1 sh[22979]: [2014-08-25 16:57:07,968: INFO/MainProcess] mingle: searching for neighbors
Aug 25 22:57:08 deis-1 sh[22979]: 2014-08-25 16:57:08 [122] [INFO] Booting worker with pid: 122
Aug 25 22:57:08 deis-1 sh[22979]: 2014-08-25 16:57:08 [123] [INFO] Booting worker with pid: 123
Aug 25 22:57:08 deis-1 sh[22979]: 2014-08-25 16:57:08 [124] [INFO] Booting worker with pid: 124
Aug 25 22:57:08 deis-1 sh[22979]: 2014-08-25 16:57:08 [125] [INFO] Booting worker with pid: 125
Aug 25 22:57:08 deis-1 sh[22979]: 2014-08-25 16:57:08 [126] [INFO] Booting worker with pid: 126
Aug 25 22:57:08 deis-1 sh[22979]: [2014-08-25 16:57:08,979: INFO/MainProcess] mingle: all alone
Aug 25 22:57:08 deis-1 sh[22979]: [2014-08-25 16:57:08,997: WARNING/MainProcess] celery@4378062f17a5 ready.
```

```console
$ deisctl journal controller
...
Aug 25 22:57:08 deis-1 sh[22979]: 2014-08-25 16:57:08 [125] [INFO] Booting worker with pid: 125
Aug 25 22:57:08 deis-1 sh[22979]: 2014-08-25 16:57:08 [126] [INFO] Booting worker with pid: 126
Aug 25 22:57:08 deis-1 sh[22979]: [2014-08-25 16:57:08,979: INFO/MainProcess] mingle: all alone
Aug 25 22:57:08 deis-1 sh[22979]: [2014-08-25 16:57:08,997: WARNING/MainProcess] celery@4378062f17a5 ready.
```

```console
$ deisctl stop controller
deis-controller@1.service: loaded
```

```console
$ deisctl start controller
deis-controller@1.service: launched
```

```console
$ deisctl scale router=3
deis-router@1.service: loaded
deis-router@2.service: loaded
deis-router@3.service: loaded

$ deisctl start router
deis-router@1.service: launched
deis-router@2.service: launched
deis-router@3.service: launched
```

## License

Copyright 2014, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
