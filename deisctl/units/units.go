package units

// Names are the base names of Deis units. Update this list when adding a new Deis unit file.
var Names = []string{
	"deis-builder",
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
}

// URL is the GitHub url where these units can be refreshed from
var URL = "https://raw.githubusercontent.com/deis/deis/"
