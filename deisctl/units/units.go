package units

// Names are the base names of Deis units. Update this list when adding a new Deis unit file.
var Names = []string{
	"deis-builder",
	"deis-cache",
	"deis-controller",
	"deis-database",
	"deis-logger",
	"deis-logspout",
	"deis-publisher",
	"deis-registry",
	"deis-router",
	"deis-store-admin",
	"deis-store-daemon",
	"deis-store-gateway",
	"deis-store-metadata",
	"deis-store-monitor",
	"deis-store-volume",
	"deis-swarm-manager",
	"deis-swarm-node",
	"deis-mesos-marathon",
	"deis-mesos-master",
	"deis-mesos-slave",
	"deis-mesos-zk@1",
}

// URL is the GitHub url where these units can be refreshed from
var URL = "https://raw.githubusercontent.com/deis/deis/%s/deisctl/units/%s.service"
