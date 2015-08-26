# Deis Logspout

A customized build of [logspout](https://github.com/gliderlabs/logspout) produced by combining *stock* logspout v3 with a custom module that adapts it for use with the [Deis](http://deis.io) open source PaaS.

**deis-logspout** is designed to run on all nodes of a Deis cluster and attach to `STDOUT` of all running [Docker](https://www.docker.com/) containers.  Logs read from these sources are forwarded to a log server such as **deis-logger**.

This Docker image is based on the official
[alpine:3.1](https://registry.hub.docker.com/_/alpine/) image.

Please add any [issues](https://github.com/deis/deis/issues) you find with this software to
the [Deis Project](https://github.com/deis/deis).  Deis maintainers will escalate issues to the upstream logspout project as needed.

## Features

The included custom logspout module implements the following components:

* `deisConfigJob`: discovers the address of the log server (e.g. **deis-logger**, which is a simple log aggregator) using [etcd](https://coreos.com/etcd/) (Deis' underlying discovery service).  After initial setup, the job continues to poll etcd at one minute intervals and will update logspout's internal configuration if the address of the log server changes.
* `deisLogAdapter`: emits log messages in the custom format expected by the **deis-logger** component.  Internally, this delegates transport to...
* `deisTransport`: a simplified version of logspout's built-in `udpTransport`.

**Note that only UDP transport is supported at this time, but support for TCP is already planned.**

## Usage

Please consult the [Makefile](Makefile) for current instructions on how to build, test, push,
install, and start **deis-logspout**.

## Environment Variables

* **DEBUG** enables debug level logging if set; this is with respect to **deis-logspout**'s own logs and *not* the logs it forwards from other containers
* **ETCD_HOST** sets the IP address of the local etcd daemon; if this is not set, etcd will not be used for log server discovery
* **ETCD_PORT** sets the TCP port on which to connect to the local etcd
  daemon (default: *4001*)

## Contributing

It is not necessary to run a Deis cluster in order to hack on **deis-logspout** locally.  All that is required is:

* A properly configured Go development environment with Go 1.4 or higher
* Docker 1.5 or higher
* Other build essentials, such a `make`

**deis-logspout** does not run in a vacuum.  To do anything useful, it must be able to attach to other containers to read their `STDOUT` and forward log messages found therein to a log server.  Additionally, since **deis-logspout** utilizes etcd for log server discovery, etcd must also be running.

In order to satisfy the dependencies noted above, **deis-logspout** comes equipped with `make` tasks for running etcd and **deis-logger** "fixtures" to better facilitate development.


### etcd

To start a single node etcd cluster:

```
make start-dev-etcd
```

This will start etcd in a daemonized Docker container.

To stop this:

```
make stop-dev-etcd
```

### deis-logger

To start **deis-logger**:

```
make run-dev-log-server
```

This launches **deis-logger** as a *non*-daemonized container since it is useful to observe this component's output during development and testing.

To stop this container, interrupt it with `ctrl-c`.

### Building and running deis-logspout

**deis-logspout** can be built and deployed in a manner consistent with all other Deis components.  Please consult the [Makefile](Makefile) for current instructions on how to build, test, push,
install, and start **deis-logspout**.

The formal build process, which produces a compact Docker image suitable for production is a relatively slow process that is inconvenient to run with any real frequency during development cycles.  To compensate, **deis-logspout** also provides an alternative build process that is optimized to run quickly during development, although it produces larger images that are not adequately compact for production use.


To build **deis-logspout** quickly:

```
make dev-build
```

Note that the first build may not be very fast, but subsequent builds will execute more quickly since they will utilized cached layers of the Docker image produced by previous builds.

To run the development-optimized build:

```
make dev-run
```

`make dev-run` will, in fact, implicitly run `make dev-build` to ensure one is always running with their latest code edits during development.


The running container can be stopped by interrupting it with `ctrl-c`.

## License

© 2015 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
