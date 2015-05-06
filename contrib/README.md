# Deis Contrib

Scripts, tools and documentation that are not part of the core
Deis system.

The contents of this directory may be vendor-specific, address a
limited audience, or be too experimental to be included in Deis' core.
This does not preclude their usefulness.

Please add any issues you find with this software to the
[Deis project](https://github.com/deis/deis/issues).

## Community Contributions

Various community members have modified Deis components, created new components, or provided tools that may be useful to Deis users. While these are not supported by the Deis team, they can be helpful in certain scenarios.

Some of these projects are listed below. This is not an exhaustive list.

Please add to this list by opening a pull request!

### Deis API Clients
* [Node.js](https://github.com/aledbf/deis-api) by [@aledbf](https://github.com/aledbf) - node.js Deis API wrapper

### Custom Deis Components
* [deis-dashboard](https://github.com/lorieri/deis-dashboard) by [@lorieri](https://github.com/lorieri) - A dashboard which summarizes requests to the Deis cluster
* [deis-docs](https://github.com/lorieri/deis-docs) by [@lorieri](https://github.com/lorieri) - Container to test Deis documentation
* [deis-netstat](https://github.com/lorieri/deis-netstat) by [@lorieri](https://github.com/lorieri) - A cluster-wide netstat tool for Deis
* [deis-proxy](https://github.com/lorieri/deis-proxy) by [@lorieri](https://github.com/lorieri) - A transparent proxy for Deis
* [deis-store-dashboard](https://github.com/aledbf/deis/tree/optional_store_dashboard) by [@aledbf](https://github.com/aledbf) - An implementation of [ceph-dash](https://github.com/Crapworks/ceph-dash) to view `deis-store` health

### CoreOS Unit Files
* [CoreOS unit files](https://github.com/ianblenke/coreos-vagrant-kitchen-sink/tree/master/cloud-init) by [@ianblenke](https://github.com/ianblenke) - Unit files to launch various services on CoreOS hosts
* [deis-backup-service](https://github.com/mozilla/deis-backup-service) by [@glogiotatidis](https://github.com/glogiotatidis) - Unit Files to automatically backup to S3 database and registry data.
* [Docker S3 Cleaner](https://github.com/myriadmobile/docker-s3-cleaner) by [@croemmich](https://github.com/croemmich) - Unit file to remove orphaned image layers from S3 backed private docker registries
* [New Relic unit for CoreOS](https://github.com/lorieri/coreos-newrelic) by [@lorieri](https://github.com/lorieri) - A global unit to launch New Relic sysmond

### Example Applications
* [Melano](https://github.com/SuaveIO/Melano) - F# "Hello World" app using the Suave framework
