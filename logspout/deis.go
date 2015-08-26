package logspout

import (
	"github.com/gliderlabs/logspout/router"
)

func init() {
	router.AdapterTransports.Register(newDeisTransport(), "deis-udp")
	router.AdapterFactories.Register(newDeisLogAdapter, "deis")
	router.Jobs.Register(newDeisConfigJob(), "deis-config-job")
}
