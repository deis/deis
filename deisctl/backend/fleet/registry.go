package fleet

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/fleet/client"
	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/pkg"
	"github.com/coreos/fleet/registry"
	"github.com/coreos/fleet/ssh"
)

// Flags used for Fleet API connectivity
var Flags = struct {
	Debug                 bool
	Version               bool
	Endpoint              string
	EtcdKeyPrefix         string
	EtcdKeyFile           string
	EtcdCertFile          string
	EtcdCAFile            string
	UseAPI                bool
	KnownHostsFile        string
	StrictHostKeyChecking bool
	Tunnel                string
	RequestTimeout        float64
	SSHTimeout            float64
}{}

const (
	oldVersionWarning = `####################################################################
WARNING: fleetctl (%s) is older than the latest registered
version of fleet found in the cluster (%s). You are strongly
recommended to upgrade fleetctl to prevent incompatibility issues.
####################################################################
`
)

func getTunnelFlag() string {
	tun := Flags.Tunnel
	if tun != "" && !strings.Contains(tun, ":") {
		tun += ":22"
	}
	return tun
}

func getChecker() *ssh.HostKeyChecker {
	if !Flags.StrictHostKeyChecking {
		return nil
	}
	keyFile := ssh.NewHostKeyFile(Flags.KnownHostsFile)
	return ssh.NewHostKeyChecker(keyFile)
}

func getFakeClient() (*registry.FakeRegistry, error) {
	return registry.NewFakeRegistry(), nil
}

func getRegistryClient() (client.API, error) {
	var dial func(string, string) (net.Conn, error)
	sshTimeout := time.Duration(Flags.SSHTimeout*1000) * time.Millisecond
	tun := getTunnelFlag()
	if tun != "" {
		sshClient, err := ssh.NewSSHClient("core", tun, getChecker(), false, sshTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed initializing SSH client: %v", err)
		}

		dial = func(network, addr string) (net.Conn, error) {
			tcpaddr, err := net.ResolveTCPAddr(network, addr)
			if err != nil {
				return nil, err
			}
			return sshClient.DialTCP(network, nil, tcpaddr)
		}
	}

	tlsConfig, err := pkg.ReadTLSConfigFiles(Flags.EtcdCAFile, Flags.EtcdCertFile, Flags.EtcdKeyFile)
	if err != nil {
		return nil, err
	}

	trans := &http.Transport{
		Dial:            dial,
		TLSClientConfig: tlsConfig,
	}

	timeout := time.Duration(Flags.RequestTimeout*1000) * time.Millisecond
	machines := []string{Flags.Endpoint}
	eClient, err := etcd.NewClient(machines, trans, timeout)
	if err != nil {
		return nil, err
	}

	reg := registry.NewEtcdRegistry(eClient, Flags.EtcdKeyPrefix)

	// if msg, ok := checkVersion(reg); !ok {
	// 	stderr(msg)
	// }

	return &client.RegistryClient{Registry: reg}, nil
}

// checkVersion makes a best-effort attempt to verify that fleetctl is at least as new as the
// latest fleet version found registered in the cluster. If any errors are encountered or fleetctl
// is >= the latest version found, it returns true. If it is < the latest found version, it returns
// false and a scary warning to the user.
// func checkVersion(reg registry.Registry) (string, bool) {
// 	fv := version.SemVersion
// 	lv, err := reg.LatestVersion()
// 	if err != nil {
// 		fmt.Printf("error attempting to check latest fleet version in Registry: %v", err)
// 	} else if lv != nil && fv.LessThan(*lv) {
// 		return fmt.Sprintf(oldVersionWarning, fv.String(), lv.String()), false
// 	}
// 	return "", true
// }
