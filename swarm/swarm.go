package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/tests/utils"
)

// EtcdCluster information about the nodes in the etcd cluster
type EtcdCluster struct {
	Members []etcd.Member `json:"members"`
}

// NodeStat information about the local node in etcd
type NodeStats struct {
	LeaderInfo struct {
		Name      string    `json:"leader"`
		Uptime    string    `json:"uptime"`
		StartTime time.Time `json:"startTime"`
	} `json:"leaderInfo"`
}

const (
	swarmpath               = "/deis/scheduler/swarm/node"
	swarmetcd               = "/deis/scheduler/swarm/host"
	etcdport                = "4001"
	timeout   time.Duration = 3 * time.Second
	ttl       time.Duration = timeout * 2
)

func run(cmd string) {
	var cmdBuf bytes.Buffer
	tmpl := template.Must(template.New("cmd").Parse(cmd))
	if err := tmpl.Execute(&cmdBuf, nil); err != nil {
		log.Fatal(err)
	}
	cmdString := cmdBuf.String()
	fmt.Println(cmdString)
	var cmdl *exec.Cmd
	cmdl = exec.Command("sh", "-c", cmdString)
	if _, _, err := utils.RunCommandWithStdoutStderr(cmdl); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("ok")
	}
}

func getleaderHost() string {
	var nodeStats NodeStats
	client := &http.Client{}
	resp, _ := client.Get("http://" + os.Getenv("HOST") + ":2379/v2/stats/self")

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &nodeStats)

	etcdLeaderID := nodeStats.LeaderInfo.Name

	var etcdCluster EtcdCluster
	resp, _ = client.Get("http://" + os.Getenv("HOST") + ":2379/v2/members")
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &etcdCluster)

	for _, node := range etcdCluster.Members {
		if node.ID == etcdLeaderID {
			u, err := url.Parse(node.ClientURLs[0])
			if err == nil {
				return u.Host
			}
		}
	}

	return ""
}

func publishService(client *etcd.Client, host string, ttl uint64) {
	for {
		setEtcd(client, swarmetcd, host, ttl)
		time.Sleep(timeout)
	}
}

func setEtcd(client *etcd.Client, key, value string, ttl uint64) {
	_, err := client.Set(key, value, ttl)
	if err != nil && !strings.Contains(err.Error(), "Key already exists") {
		log.Println(err)
	}
}

func main() {
	etcdproto := "etcd://" + getleaderHost() + swarmpath
	etcdhost := os.Getenv("HOST")
	addr := "--addr=" + etcdhost + ":2375"
	client := etcd.NewClient([]string{"http://" + etcdhost + ":" + etcdport})
	switch os.Args[1] {
	case "join":
		run("./deis-swarm join " + addr + " " + etcdproto)
	case "manage":
		go publishService(client, etcdhost, uint64(ttl.Seconds()))
		run("./deis-swarm manage " + etcdproto)
	}
}
