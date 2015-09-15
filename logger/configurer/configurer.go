package configurer

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/logger/drain"
	"github.com/deis/deis/logger/storage"
	"github.com/deis/deis/logger/syslogish"
)

// Exported so it can be set by an external agent-- namely main.go, which does some flag parsing.
var DefaultDrainURI string

// Configurer takes responsibility for dynamically reconfiguring a syslogish.Server based on
// changes in etcd.
type Configurer struct {
	etcdClient                *etcd.Client
	etcdPath                  string
	ticker                    *time.Ticker
	syslogishServer           *syslogish.Server
	running                   bool
	currentStorageAdapterType string
	currentDrainURL           string
}

// NewConfigurer returns a pointer to a new Configurer instance.
func NewConfigurer(etcdHost string, etcdPort int, etcdPath string, configInterval int,
	syslogishServer *syslogish.Server) (*Configurer, error) {
	etcdClient := etcd.NewClient([]string{fmt.Sprintf("http://%s:%d", etcdHost, etcdPort)})
	ticker := time.NewTicker(time.Duration(configInterval) * time.Second)
	configurer := &Configurer{
		etcdClient:      etcdClient,
		etcdPath:        etcdPath,
		syslogishServer: syslogishServer,
		ticker:          ticker,
	}

	// Support legacy behavior that allows default drain uri to be specified using a drain-uri flag
	if _, err := etcdClient.Get(etcdPath+"/drain", false, false); err != nil {
		etcdErr, ok := err.(*etcd.EtcdError)
		// Error code 100 is key not found
		if ok && etcdErr.ErrorCode == 100 {
			configurer.setEtcd("/drain", DefaultDrainURI)
		} else {
			log.Println(err)
		}
	}

	return configurer, nil
}

// Start begins the configurer's main loop.
func (c *Configurer) Start() {
	// Should only ever be called once
	if !c.running {
		c.running = true
		go c.configure()
		log.Println("configurer running")
	}
}

func (c *Configurer) configure() {
	for {
		<-c.ticker.C
		c.manageStorageAdapter()
		c.manageDrain()
	}
}

func (c *Configurer) manageStorageAdapter() {
	newStorageAdapterType, err := c.getEtcd("/storageAdapterType", "file")
	if err != nil {
		log.Println("configurer: Error retrieving storage adapter type from etcd.  Skipping.", err)
		return
	}
	if newStorageAdapterType == c.currentStorageAdapterType {
		return
	}
	newStorageAdapter, err := storage.NewAdapter(newStorageAdapterType)
	if err != nil {
		log.Println("configurer: Error creating new storage adapter.  Skipping.", err)
		return
	}
	c.syslogishServer.SetStorageAdapter(newStorageAdapter)
	c.currentStorageAdapterType = newStorageAdapterType
	log.Printf("configurer: Activated new storage adapter: %s", newStorageAdapterType)
}

func (c *Configurer) manageDrain() {
	newDrainURL, err := c.getEtcd("/drain", "")
	if err != nil {
		log.Println("configurer: Error retrieving drain URL from etcd.  Skipping.", err)
		return
	}
	if newDrainURL == c.currentDrainURL {
		return
	}
	newDrain, err := drain.NewDrain(newDrainURL)
	if err != nil {
		log.Println("configurer: Error creating new drain.  Skipping.", err)
		return
	}
	c.syslogishServer.SetDrain(newDrain)
	c.currentDrainURL = newDrainURL
	if newDrainURL == "" {
		log.Println("configurer: Deactivated drain")
	} else {
		log.Printf("configurer: Activated new drain: %s", newDrainURL)
	}
}

func (c *Configurer) getEtcd(key string, defaultValue string) (string, error) {
	resp, err := c.etcdClient.Get(fmt.Sprintf("%s%s", c.etcdPath, key), false, false)
	if err != nil {
		etcdErr, ok := err.(*etcd.EtcdError)
		// Error code 100 is key not found
		if ok && etcdErr.ErrorCode == 100 {
			return defaultValue, nil
		}
		return "", err
	}
	return resp.Node.Value, nil
}

func (c *Configurer) setEtcd(key string, value string) {
	_, err := c.etcdClient.Set(fmt.Sprintf("%s%s", c.etcdPath, key), value, 0)
	if err != nil {
		etcdErr, ok := err.(*etcd.EtcdError)
		// Error code 105 is key already exists
		if !ok || etcdErr.ErrorCode != 105 {
			log.Println(err)
		}
	}
}
