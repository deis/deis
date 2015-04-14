# Mock S3 storage

The objective is to provide an S3-compatible service for tests, so a Ceph cluster
is not required. This component uses [mock-s3](https://github.com/jserver/mock-s3).

## Usage:

```
docker run -p 8888:8888 -e HOST=$COREOS_PRIVATE_IPV4 -v <local directory>:/app/storage deis/mock-store
```

*The use of a local directory `(-v <local directory>)` is optional*


`mock-s3` does not requires an `ACCESS_KEY` and `SECRET_KEY` (there is no concept of permissions), but this
component will generate both to keep compatibility with `deis-store-gateway`.

## Containers

The mock store component is composed of one container:

* [mock-store](https://index.docker.io/u/deis/mock-store/) - the blob store gateway,
offering a S3-compatible bucket APIs using the local filesystem as storage.

## Usage

Please consult the [Makefile](Makefile) for current instructions on how to build, test, push,
install, and start **deis/mock-store**.

Note that changes to **deis/mock-store** will *not* be built automatically by the test suite.
Run `make mock-store` from the tests/ directory to update the Docker image used by tests.

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
