package etcd

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/fleet/pkg"
	"github.com/coreos/fleet/ssh"
	etcdlib "github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/deisctl/backend/fleet"
	"github.com/deis/deis/deisctl/config/model"
)

// ConfigBackend is an etcd-based implementation of the config.Backend interface
type ConfigBackend struct {
	etcdlib *etcdlib.Client
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

// Get a value by key from etcd
func (cb *ConfigBackend) Get(key string) (string, error) {
	sort, recursive := true, false
	resp, err := cb.etcdlib.Get(key, sort, recursive)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// GetWithDefault gets a value by key from etcd and return a default value if
// not found
func (cb *ConfigBackend) GetWithDefault(key string, defaultValue string) (string, error) {
	sort, recursive := true, false
	resp, err := cb.etcdlib.Get(key, sort, recursive)
	if err != nil {
		etcdErr, ok := err.(*etcdlib.EtcdError)
		if ok && etcdErr.ErrorCode == 100 {
			return defaultValue, nil
		}
		return "", err
	}
	return resp.Node.Value, nil
}

func singleNodeToConfigNode(node *etcdlib.Node) *model.ConfigNode {
	key := model.ConfigNode{
		Key:        node.Key,
		Expiration: node.Expiration,
	}

	if node.Dir != true && node.Key != "" {
		key.Value = node.Value
	}

	return &key
}

func traverseNode(node *etcdlib.Node) []*model.ConfigNode {
	var serviceKeys []*model.ConfigNode

	if len(node.Nodes) > 0 {
		for _, nodeChild := range node.Nodes {
			serviceKeys = append(serviceKeys, traverseNode(nodeChild)...)
		}
	} else {
		key := singleNodeToConfigNode(node)
		if key.Key != "" {
			serviceKeys = append(serviceKeys, key)
		}
	}

	return serviceKeys
}

// GetRecursive returns a slice of all key/value pairs "under" a specified key
// in etcd
func (cb *ConfigBackend) GetRecursive(key string) ([]*model.ConfigNode, error) {
	resp, err := cb.etcdlib.Get(key, true, true)
	if err != nil {
		return nil, err
	}

	nodes := traverseNode(resp.Node)
	return nodes, nil
}

// Delete a key/value pair by key from etcd
func (cb *ConfigBackend) Delete(key string) error {
	_, err := cb.etcdlib.Delete(key, false)
	return err
}

// Set a value for the specified key in etcd
func (cb *ConfigBackend) Set(key string, value string) (string, error) {
	resp, err := cb.etcdlib.Set(key, value, 0) // don't use TTLs
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// SetWithTTL sets a value for the specified key in etcd-- with a time to live
func (cb *ConfigBackend) SetWithTTL(key string, value string, ttl uint64) (string, error) {
	resp, err := cb.etcdlib.Update(key, value, ttl)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// NewConfigBackend returns this etcd-based implementation of the config.Backend
// interface
func NewConfigBackend() (*ConfigBackend, error) {
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

	c := etcdlib.NewClient(machines)
	c.SetDialTimeout(timeout)

	// use custom transport with SSH tunnel capability
	c.SetTransport(&trans)

	return &ConfigBackend{etcdlib: c}, nil
}
