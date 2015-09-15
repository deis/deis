package publisher

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

// Publisher takes responsibility for regularly updating etcd with the host and port where this
// logger component is running.  This permits other components to discover it.
type Publisher struct {
	etcdClient *etcd.Client
	etcdPath   string
	publishTTL uint64
	logHost    string
	logPort    int
	ticker     *time.Ticker
	running    bool
}

// NewPublisher returns a pointer to a new Publisher instance.
func NewPublisher(etcdHost string, etcdPort int, etcdPath string, publishInterval int,
	publishTTL int, logHost string, logPort int) (*Publisher, error) {
	etcdClient := etcd.NewClient([]string{fmt.Sprintf("http://%s:%d", etcdHost, etcdPort)})
	ticker := time.NewTicker(time.Duration(publishInterval) * time.Second)
	return &Publisher{
		etcdClient: etcdClient,
		etcdPath:   etcdPath,
		publishTTL: uint64(time.Duration(publishTTL) * time.Second),
		logHost:    logHost,
		logPort:    logPort,
		ticker:     ticker,
	}, nil
}

// Start begins the publisher's main loop.
func (p *Publisher) Start() {
	// Should only ever be called once
	if !p.running {
		p.running = true
		go p.publish()
		log.Println("publisher running")
	}
}

func (p *Publisher) publish() {
	for {
		<-p.ticker.C
		p.setEtcd("/host", p.logHost)
		p.setEtcd("/port", strconv.Itoa(p.logPort))
	}
}

func (p *Publisher) setEtcd(key string, value string) {
	_, err := p.etcdClient.Set(fmt.Sprintf("%s%s", p.etcdPath, key), value, p.publishTTL)
	if err != nil {
		etcdErr, ok := err.(*etcd.EtcdError)
		// Error code 105 is key already exists
		if !ok || etcdErr.ErrorCode != 105 {
			log.Println(err)
		}
	}
}
