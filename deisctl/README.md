# Deis Control Utility

`deisctl` is a command-line utility used to operate a Deis cluster.

## Build & Install on Deis Hosts

Requirements:

* GNU make
* Docker client

```
$ make package
docker build -t deis/deisctl .
Uploading context 9.698 MB
Uploading context
Step 0 : FROM deis/base
...
package placed in ./package/deisctl.tar.gz
$ tar tvfz ./package/deisctl.tar.gz
...<unit files>...
drwxr-xr-x  0 root   root        0 Jul  1 11:17 ./opt/
drwxr-xr-x  0 root   root        0 Jul  1 11:17 ./opt/bin/
-rwxr-xr-x  0 root   root  8789856 Jul  1 11:17 ./opt/bin/deisctl
```

Copy the tarball to the CoreOS host and untar it from the root.
For example, with a local Vagrant environment:

```
$ vagrant rsync
==> deis-1: Rsyncing folder: /Users/gabriel/workspace/src/github.com/deis/deis/ => /home/core/share
$ vagrant ssh -c "sudo tar -C / -xvf /home/core/share/deisctl/package/deisctl.tar.gz"
./
./var/
./var/lib/
./var/lib/deis/
./var/lib/deis/units/
./var/lib/deis/units/builder/
./var/lib/deis/units/builder/deis-builder.service
./var/lib/deis/units/cache/
./var/lib/deis/units/cache/deis-cache.service
./var/lib/deis/units/controller/
./var/lib/deis/units/controller/deis-controller.service
./var/lib/deis/units/database/
./var/lib/deis/units/database/deis-database.service
./var/lib/deis/units/logger/
./var/lib/deis/units/logger/deis-logger.service
./var/lib/deis/units/registry/
./var/lib/deis/units/registry/deis-registry.service
./var/lib/deis/units/router/
./var/lib/deis/units/router/deis-router.service
./opt/
./opt/bin/
./opt/bin/deisctl

$ vagrant ssh -c deisctl
Usage:
  deisctl <command> [<target>...] [options]
```

## Build & Install on your Local Workstation

Requirements:

* GNU make
* Go 1.2+ runtime with $GOPATH/bin in your shell path
* Godep for managing external dependencies (install with `go get github.com/tools/godep`)

```
$ make install
godep go install ./...
binaries placed in $GOPATH/bin

$
```

Note `deisctl` uses the `FLEETCTL_TUNNEL` environment variable for remote connectivity.
For example, to connect to a local Vagrant host from your workstation:

```
$ export FLEETCTL_TUNNEL=172.17.8.100
$ deisctl list
UNIT				STATE		LOAD	ACTIVE	SUB	DESC		MACHINE
deis-builder.1.service		launched	loaded	active	running	deis-builder	2f603f6e.../172.17.8.100
deis-cache.1.service		launched	loaded	active	running	deis-cache	2f603f6e.../172.17.8.100
deis-controller.1.service	launched	loaded	active	running	deis-controller	2f603f6e.../172.17.8.100
deis-database.1.service		launched	loaded	active	running	deis-database	2f603f6e.../172.17.8.100
deis-logger.1.service		launched	loaded	active	running	deis-logger	2f603f6e.../172.17.8.100
deis-registry.1.service		launched	loaded	active	running	deis-registry	2f603f6e.../172.17.8.100
deis-router.1.service		launched	loaded	active	running	deis-router	2f603f6e.../172.17.8.100
```

## Usage

To install all Deis units, use `deisctl install`.  Other commands can be found via `deisctl help`:

```
$ deisctl help
Deis Control Utility

Usage:
  deisctl <command> [<target>...] [options]

Example Commands:

  deisctl install
  deisctl uninstall
  deisctl list
  deisctl scale router=2
  deisctl start router.2
  deisctl stop router builder
  deisctl status controller

Options:
  --debug                     print debug information to stderr
  --endpoint=<url>            etcd endpoint for fleet [default: http://127.0.0.1:4001]
  --etcd-key-prefix=<path>    keyspace for fleet data in etcd [default: /_coreos.com/fleet/]
  --known-hosts-file=<path>   file used to store remote machine fingerprints [default: ~/.fleetctl/known_hosts]
  --strict-host-key-checking  verify SSH host keys [default: true]
  --tunnel=<host>             establish an SSH tunnel for communication with fleet and etcd
  --verbosity=<level>         log at a specified level of verbosity to stderr [default: 0]
```

## License

Copyright 2014, OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
