package etcd

import (
	"errors"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-etcd/etcd"
	logger "github.com/deis/deis/mesos/pkg/log"
	etcdlock "github.com/leeor/etcd-sync"
)

// Client etcd client
type Client struct {
	client *etcd.Client
	lock   *etcdlock.EtcdMutex
}

// Error etcd error
type Error struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
	Cause     string `json:"cause,omitempty"`
	Index     uint64 `json:"index"`
}

var log = logger.New()

// NewClient create a etcd client using the given machine list
func NewClient(machines []string) *Client {
	log.Debugf("connecting to %v etcd server/s", machines)
	return &Client{etcd.NewClient(machines), nil}
}

// SetDefault sets the value of a key without expiration
func SetDefault(client *Client, key, value string) {
	Create(client, key, value, 0)
}

// Mkdir creates a directory only if does not exists
func Mkdir(c *Client, path string) {
	_, err := c.client.CreateDir(path, 0)
	if err != nil {
		log.Debug(err)
	}
}

// WaitForKeys wait for the required keys up to the timeout or forever if is nil
func WaitForKeys(c *Client, keys []string, ttl time.Duration) error {
	start := time.Now()
	wait := true

	for {
		for _, key := range keys {
			_, err := c.client.Get(key, false, false)
			if err != nil {
				log.Debugf("key \"%s\" error %v", key, err)
				wait = true
			}
		}

		if !wait {
			return nil
		}

		log.Debug("waiting for missing etcd keys...")
		time.Sleep(1 * time.Second)
		wait = false

		if time.Since(start) > ttl {
			return errors.New("maximum ttl reached. aborting")
		}
	}
}

// Get returns the value inside a key or an empty string
func Get(c *Client, key string) string {
	result, err := c.client.Get(key, false, false)
	if err != nil {
		log.Debugf("%v", err)
		return ""
	}

	return result.Node.Value
}

// GetList returns the list of elements inside a key or an empty list
func GetList(c *Client, key string) []string {
	values, err := c.client.Get(key, true, false)
	if err != nil {
		log.Debugf("getlist %v", err)
		return []string{}
	}

	result := []string{}
	for _, node := range values.Node.Nodes {
		result = append(result, path.Base(node.Key))
	}

	log.Debugf("getlist %s -> %v", key, result)
	return result
}

// Set sets the value of a key.
// If the ttl is bigger than 0 it will expire after the specified time
func Set(c *Client, key, value string, ttl uint64) {
	log.Debugf("set %s -> %s", key, value)
	_, err := c.client.Set(key, value, ttl)
	if err != nil {
		log.Debugf("%v", err)
	}
}

// Create set the value of a key only if it does not exits
func Create(c *Client, key, value string, ttl uint64) {
	log.Debugf("create %s -> %s", key, value)
	_, err := c.client.Create(key, value, ttl)
	if err != nil {
		log.Debugf("%v", err)
	}
}

// PublishService publish a service to etcd periodically
func PublishService(
	client *Client,
	etcdPath string,
	host string,
	externalPort int,
	ttl uint64,
	timeout time.Duration) {

	for {
		Set(client, etcdPath+"/host", host, ttl)
		Set(client, etcdPath+"/port", strconv.Itoa(externalPort), ttl)
		time.Sleep(timeout)
	}
}

func convertEtcdError(err error) *Error {
	etcdError := err.(*etcd.EtcdError)
	return &Error{
		ErrorCode: etcdError.ErrorCode,
		Message:   etcdError.Message,
		Cause:     etcdError.Cause,
		Index:     etcdError.Index,
	}
}

// GetHTTPEtcdUrls returns an array of urls that contains at least one host
func GetHTTPEtcdUrls(host, etcdPeers string) []string {
	if etcdPeers != "127.0.0.1:4001" {
		hosts := strings.Split(etcdPeers, ",")
		result := []string{}
		for _, _host := range hosts {
			result = append(result, "http://"+_host+":4001")
		}
		return result
	}

	return []string{"http://" + host}
}
