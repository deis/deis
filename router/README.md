# Deis Router

An nginx proxy for use in the [Deis](http://deis.io) open source PaaS.

[![image](https://d207aa93qlcgug.cloudfront.net/img/icons/framed-icon-checked-repository.svg)](https://index.docker.io/u/deis/router/)

[**Trusted Build**](https://index.docker.io/u/deis/router/)

This Docker image is based on the trusted build
[deis/base](https://index.docker.io/u/deis/base/), which itself is based
on the official [ubuntu:12.04](https://index.docker.io/_/ubuntu/) image.

Please add any issues you find with this software to the
[Deis project](https://github.com/deis/deis/issues).

## Usage

* `make build` builds the *deis/router* image inside a vagrant VM
* `make run` installs and starts *deis/router*, then displays log
  output from the container

## Environment Variables

* **DEBUG** enables verbose output if set
* **ETCD_PORT** sets the TCP port on which to connect to the local etcd
  daemon (default: *4001*)
* **ETCD_PATH** sets the etcd directory where the router announces
  its configuration (default: */deis/router*)
* **ETCD_TTL** sets the time-to-live before etcd purges a configuration
  value, in seconds (default: *10*)
* **PORT** sets the TCP port on which the router listens (default: *80*)


## License

© 2014 OpDemand LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
