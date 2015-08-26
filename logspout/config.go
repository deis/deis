package logspout

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/gliderlabs/logspout/router"
)

type deisConfigJob struct {
	etcd *etcd.Client
}

func newDeisConfigJob() *deisConfigJob {
	return &deisConfigJob{}
}

func (j *deisConfigJob) Name() string {
	return "deis-config"
}

func (j *deisConfigJob) Setup() error {
	if etcdHost := os.Getenv("ETCD_HOST"); etcdHost != "" {
		etcdPort := os.Getenv("ETCD_PORT")
		if etcdPort == "" {
			etcdPort = "4001"
		}
		connectionString := []string{"http://" + etcdHost + ":" + etcdPort}
		fmt.Println("# deis-config: etcd: " + connectionString[0])
		j.etcd = etcd.NewClient(connectionString)
		j.etcd.SetDialTimeout(3 * time.Second)
	} else {
		fmt.Println("# deis-config: etcd: no connection details provided -- NOT USING")
	}
	return nil
}

func (j *deisConfigJob) Run() error {
	var currentRoute *router.Route
	for {
		if j.etcd != nil {
			hostResp, err := j.etcd.Get("/deis/logs/host", false, false)
			if err != nil {
				log.Println("deis-config: etcd:", err)
			} else {
				portResp, err := j.etcd.Get("/deis/logs/port", false, false)
				if err != nil {
					log.Println("deis-config: etcd:", err)
				} else {
					host := fmt.Sprintf("%s:%s", hostResp.Node.Value, portResp.Node.Value)
					if currentRoute == nil || host != currentRoute.Address {
						if currentRoute != nil {
							router.Routes.Remove(currentRoute.ID)
							log.Println("deis-config: removed route to", currentRoute.Address)
						}
						currentRoute = &router.Route{Address: host, Adapter: "deis", Options: make(map[string]string)}
						if err := router.Routes.Add(currentRoute); err != nil {
							log.Println("deis-config:", err)
						} else {
							log.Println("deis-config: added route to", currentRoute.Address)
						}
					}
				}
			}
		}
		time.Sleep(60 * time.Second)
	}
	return nil
}
