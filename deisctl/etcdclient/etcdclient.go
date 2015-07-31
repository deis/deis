package etcdclient

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/fleet/pkg"
	"github.com/coreos/fleet/ssh"
	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/deisctl/backend/fleet"
)

// ServiceKey represents running Deis services
type ServiceKey struct {
	Key        string     `json:"key"`
	Value      string     `json:"value,omitempty"`
	Expiration *time.Time `json:"expiration,omitempty"`
	TTL        int64      `json:"ttl,omitempty"`
}

func getTunnelFlag() string {
	tun := fleet.Flags.Tunnel
	if tun != "" && !strings.Contains(tun, ":") {
		tun += ":22"
	}
	return tun
}

func getChecker() *ssh.HostKeyChecker {
	if !fleet.Flags.StrictHostKeyChecking {
		return nil
	}
	keyFile := ssh.NewHostKeyFile(fleet.Flags.KnownHostsFile)
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

func singleNodeToServiceKey(node *etcd.Node) *ServiceKey {
	key := ServiceKey{
		Key:        node.Key,
		Expiration: node.Expiration,
	}

	if node.Dir != true && node.Key != "" {
		key.Value = node.Value
	}

	return &key
}

func traverseNode(node *etcd.Node) []*ServiceKey {
	var serviceKeys []*ServiceKey

	if len(node.Nodes) > 0 {
		for _, nodeChild := range node.Nodes {
			serviceKeys = append(serviceKeys, traverseNode(nodeChild)...)
		}
	} else {
		key := singleNodeToServiceKey(node)
		if key.Key != "" {
			serviceKeys = append(serviceKeys, key)
		}
	}

	return serviceKeys
}

func (c *etcdClient) GetRecursive(key string) ([]*ServiceKey, error) {
	resp, err := c.etcd.Get(key, true, true)
	if err != nil {
		return nil, err
	}

	nodes := traverseNode(resp.Node)
	return nodes, nil
}

func (c *etcdClient) Delete(key string) error {
	_, err := c.etcd.Delete(key, false)
	return err
}

func (c *etcdClient) Set(key string, value string) (string, error) {
	resp, err := c.etcd.Set(key, value, 0) // don't use TTLs
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func (c *etcdClient) Update(key string, value string, ttl uint64) (string, error) {
	resp, err := c.etcd.Update(key, value, ttl)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// GetEtcdClient returns a valid etcd client, either locally or via SSH
func GetEtcdClient() (*etcdClient, error) {
	var dial func(string, string) (net.Conn, error)
	sshTimeout := time.Duration(fleet.Flags.SSHTimeout*1000) * time.Millisecond
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

	tlsConfig, err := pkg.ReadTLSConfigFiles(fleet.Flags.EtcdCAFile,
		fleet.Flags.EtcdCertFile, fleet.Flags.EtcdKeyFile)
	if err != nil {
		return nil, err
	}

	trans := http.Transport{
		Dial:            dial,
		TLSClientConfig: tlsConfig,
	}

	timeout := time.Duration(fleet.Flags.RequestTimeout*1000) * time.Millisecond
	machines := []string{fleet.Flags.Endpoint}

	c := etcd.NewClient(machines)
	c.SetDialTimeout(timeout)

	// use custom transport with SSH tunnel capability
	c.SetTransport(&trans)

	return &etcdClient{etcd: c}, nil
}
