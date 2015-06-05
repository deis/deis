# Deis Router

An nginx proxy for use in the [Deis](http://deis.io) open source PaaS.

This Docker image is based on the official
[alpine:3.2](https://registry.hub.docker.com/_/alpine/) image.

Please add any [issues](https://github.com/deis/deis/issues) you find with this software to
the [Deis Project](https://github.com/deis/deis).

## Usage

Please consult the [Makefile](Makefile) for current instructions on how to build, test, push,
install, and start **deis/router**.

## Environment Variables

* **DEBUG** enables verbose output if set
* **ETCD_PORT** sets the TCP port on which to connect to the local etcd
  daemon (default: *4001*)
* **ETCD_PATH** sets the etcd directory where the router announces
  its configuration (default: */deis/router*)
* **ETCD_TTL** sets the time-to-live before etcd purges a configuration
  value, in seconds (default: *10*)
* **PORT** sets the TCP port on which the router listens (default: *80*)


## Firewall

[Shellshock](https://shellshocker.net) exposed that some apps (mostly CGI based) inside a web server can be exploited, allowing the arbitrary execution of commands.

To reduce the contact surface of this attack and others (like SQL injection and cross site scripting), it's possible to enable the naxsi firewall (which is disabled by default). [**NAXSI**](https://github.com/nbs-system/naxsi) is an open-source, high performance, low rules maintenance WAF for NGINX.
The rules included are from this project [doxi-rules](https://bitbucket.org/lazy_dogtown/doxi-rules)

Only these modules are enabled:

|File|Purpose|
|----|-------|
|web_app.rules       |detect exploit/misuse-attempts againts web-applications
|web_server.rules    |generic rules to protect a webserver from misconfiguration and known mistakes / exploit-vectors
|active-mode.rules   |rules to configure active-mode (block)
|naxsi_core          |core naxsi rules

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
