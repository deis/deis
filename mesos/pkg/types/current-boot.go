package types

import (
	"net"
	"time"

	"github.com/deis/deis/mesos/pkg/etcd"
)

// CurrentBoot information about the boot
// process related to the component
type CurrentBoot struct {
	ConfdNodes []string
	EtcdClient *etcd.Client
	EtcdPath   string
	EtcdPort   int
	EtcdPeers  string
	EtcdURL    []string
	Host       net.IP
	Port       int
	Timeout    time.Duration
	TTL        time.Duration
}
