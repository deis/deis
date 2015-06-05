# Deis Database

A PostgreSQL database for use in the [Deis](http://deis.io) open source PaaS.

This Docker image is based on the official
[alpine:3.2](https://registry.hub.docker.com/_/alpine/) image.

Please add any [issues](https://github.com/deis/deis/issues) you find with this software to
the [Deis Project](https://github.com/deis/deis).

## Usage

Please consult the [Makefile](Makefile) for current instructions on how to build, test, push,
install, and start **deis/database**.

## Environment Variables

* **BUCKET_NAME** store component bucket used for database WAL logs and backups (default: db_wal)
* **DEBUG** enables verbose output if set
* **ETCD_PORT** sets the TCP port on which to connect to the local etcd
  daemon (default: *4001*)
* **ETCD_PATH** sets the etcd directory where the database announces
  its configuration (default: */deis/database*)
* **ETCD_TTL** sets the time-to-live before etcd purges a configuration
  value, in seconds (default: *10*)
* **PG_ADMIN_USER** sets the database admin user name (default: *postgres*)
* **PG_ADMIN_PASS** sets the database admin user password
  (default: *changeme123*)
* **PG_CONFIG** sets the PostgreSQL configuration file location
  (default: */etc/postgresql/9.3/main/postgresql.conf*)
* **PG_LISTEN** sets the addresses on which the database will listen
  (default: *)
* **PG_USER_NAME** sets the database user name for Deis (default: *deis*)
* **PG_USER_PASS** sets the database user password for Deis
  (default: *changeme123*)
* **PG_USER_DB** sets the database name used by Deis (default: *deis*)
* **PORT** sets the TCP port on which the database listens (default: *5432*)

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
