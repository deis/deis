# Deis Registry

A Docker image registry for use in the [Deis](http://deis.io) open source PaaS.

This Docker image is based on the official
[alpine:3.2](https://registry.hub.docker.com/_/alpine/) image.

Please add any [issues](https://github.com/deis/deis/issues) you find with this software to
the [Deis Project](https://github.com/deis/deis).

## Usage

Please consult the [Makefile](Makefile) for current instructions on how to build, test, push,
install, and start **deis/registry**.

## Environment Variables

* **DEBUG** enables verbose output if set
* **ETCD_PORT** sets the TCP port on which to connect to the local etcd
  daemon (default: *4001*)
* **ETCD_PATH** sets the etcd directory where the registry announces
  its configuration (default: */deis/registry*)
* **ETCD_TTL** sets the time-to-live before etcd purges a configuration
  value, in seconds (default: *10*)
* **PORT** sets the TCP port on which the registry listens (default: *5000*)
* **REGISTRY_PORT** set the TCP port on which registry listens
  (default: *5000*)
* **REGISTRY_SECRET_KEY** set the registry's secret key (default: randomized)

## License

© 2014 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
