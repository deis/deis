package types

import (
	"net"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

// CurrentBoot information about the boot
// process related to the component
type CurrentBoot struct {
	EtcdClient *etcd.Client
	EtcdPath   string
	EtcdPort   string
	Host       net.IP
	Port       int
	Timeout    time.Duration
	TTL        time.Duration
}
