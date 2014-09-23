# Deis Store

A backing store built on [Ceph](http://ceph.com) for use in the [Deis](http://deis.io) open
source PaaS.

The bin/boot scripts and Dockerfiles were inspired by Seán C. McCord's [docker-ceph](https://github.com/Ulexus/docker-ceph) repository.

Please add any issues you find with this software to the
[Deis project](https://github.com/deis/deis/issues).

## Containers

The store component is comprised of four containers:

* [store-daemon](https://index.docker.io/u/deis/store-daemon/) - the daemon which serves data (in Ceph, this is an object store daemon, or OSD)
* [store-monitor](https://index.docker.io/u/deis/store-monitor/) - the service responsible for keeping track of the cluster state (this is also called a monitor in Ceph)
* [store-gateway](https://index.docker.io/u/deis/store-gateway/) - the blob store gateway, offering Swift and S3-compatible bucket APIs
* [store-data](systemd/deis-store-data@.service.template) - a container which exposes volumes to be used by the other containers

These are all based upon the [store-base](https://index.docker.io/u/deis/store-base/) image, which is a Docker container that preinstalls Ceph.

## Usage

The Makefile in this directory knows how to build the containers. The Deis root Makefile knows how to install and start them:

```console
$ make -C store build
$ make install-stores start-stores
```

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
