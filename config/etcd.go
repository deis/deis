package config

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	fleetEtcd "github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/ssh"
	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deisctl/client"
)

func getTunnelFlag() string {
	tun := client.Flags.Tunnel
	if tun != "" && !strings.Contains(tun, ":") {
		tun += ":22"
	}
	return tun
}

func getChecker() *ssh.HostKeyChecker {
	if !client.Flags.StrictHostKeyChecking {
		return nil
	}
	keyFile := ssh.NewHostKeyFile(client.Flags.KnownHostsFile)
	return ssh.NewHostKeyChecker(keyFile)
}

type etcdClient struct {
	etcd *etcd.Client
}

func (c *etcdClient) Get(key string) (string, error) {
	sort, recursive := true, false
	resp, err := c.etcd.Get(key, sort, recursive)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func (c *etcdClient) Set(key string, value string) (string, error) {
	resp, err := c.etcd.Set(key, value, 0) // don't use TTLs
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func getEtcdClient() (*etcdClient, error) {
	var dial func(string, string) (net.Conn, error)
	tun := getTunnelFlag()
	if tun != "" {
		sshClient, err := ssh.NewSSHClient("core", tun, getChecker(), false)
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

	tlsConfig, err := fleetEtcd.ReadTLSConfigFiles(client.Flags.EtcdCAFile,
		client.Flags.EtcdCertFile, client.Flags.EtcdKeyFile)
	if err != nil {
		return nil, err
	}

	trans := http.Transport{
		Dial:            dial,
		TLSClientConfig: tlsConfig,
	}

	timeout := time.Duration(client.Flags.RequestTimeout*1000) * time.Millisecond
	machines := []string{client.Flags.Endpoint}

	c := etcd.NewClient(machines)
	c.SetDialTimeout(timeout)

	// use custom transport with SSH tunnel capability
	c.SetTransport(&trans)

	return &etcdClient{etcd: c}, nil
}
