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

### Backup Tools
* [Deis Backup and Restore](https://github.com/myriadmobile/deis-backup-restore) by [@croemmich](https://github.com/croemmich) - Container to backup and restore etcd, database, registry, and logs to/from any S3 compatible store.
* [deis-backup-service](https://github.com/mozilla/deis-backup-service) by [@glogiotatidis](https://github.com/glogiotatidis) - Container to backup the database and registry. Uses s3cmd behind the scenes and supports data encryption.

### Deis API Clients
* [Node.js](https://github.com/aledbf/deis-api) by [@aledbf](https://github.com/aledbf) - node.js Deis API wrapper

### Custom Deis Components
* [deis-dashboard](https://github.com/lorieri/deis-dashboard) by [@lorieri](https://github.com/lorieri) - A dashboard which summarizes requests to the Deis cluster
* [deis-docs](https://github.com/lorieri/deis-docs) by [@lorieri](https://github.com/lorieri) - Container to test Deis documentation
* [deis-netstat](https://github.com/lorieri/deis-netstat) by [@lorieri](https://github.com/lorieri) - A cluster-wide netstat tool for Deis
* [deis-proxy](https://github.com/lorieri/deis-proxy) by [@lorieri](https://github.com/lorieri) - A transparent proxy for Deis
* [deis-store-dashboard](https://github.com/aledbf/deis/tree/optional_store_dashboard) by [@aledbf](https://github.com/aledbf) - An implementation of [ceph-dash](https://github.com/Crapworks/ceph-dash) to view `deis-store` health
* [deis-phppgadmin](https://github.com/HeheCloud/deis-phppgadmin) by [hehecloud](https://github.com/HeheCloud) - An addon (database dashboard) for deis-database (phpPgAdmin)


### CoreOS Unit Files
* [CoreOS unit files](https://github.com/ianblenke/coreos-vagrant-kitchen-sink/tree/master/cloud-init) by [@ianblenke](https://github.com/ianblenke) - Unit files to launch various services on CoreOS hosts
* [Docker S3 Cleaner](https://github.com/myriadmobile/docker-s3-cleaner) by [@croemmich](https://github.com/croemmich) - Unit file to remove orphaned image layers from S3 backed private docker registries
* [New Relic unit for CoreOS](https://github.com/lorieri/coreos-newrelic) by [@lorieri](https://github.com/lorieri) - A global unit to launch New Relic sysmond
* [Sematext Docker Agent for CoreOS](https://github.com/sematext/sematext-agent-docker/blob/master/coreos/sematext-agent.service) by [@sematext](https://github.com/sematext) - A global unit to launch the agent for [SPM Performance Monitoring, Anomaly Detection and Alerting](http://sematext.com/spm/integrations/docker-monitoring.html) 
* [Forwarding systemd journal to Logsene](https://github.com/sematext/sematext-agent-docker/blob/master/coreos/logsene.service) by [@sematext](https://github.com/sematext) - A global unit to forward systemd journal via SSL/TLS. Note: The IP address of the CoreOS host needs to be authorized in Logsene. [Logsene ­Log Management & Analytics](http://www.sematext.com/logsene/) 

### Example Applications
* [Melano](https://github.com/SuaveIO/Melano) - F# "Hello World" app using the Suave framework
