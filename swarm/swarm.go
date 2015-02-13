package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/tests/utils"
)

//swarm join --addr=<node_ip:2375> etcd://<etcd_ip>/<path>
type Response struct {
	ClientURL string `json:"clientURL"`
	Name      string `json:"name"`
	PeerURL   string `json:"peerURL"`
	State     string `json:"state"`
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
	var response []Response
	var host string
	client := &http.Client{}
	resp, _ := client.Get("http://" + os.Getenv("HOST") + ":7001/v2/admin/machines")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// body, _ := ioutil.ReadFile("test.json")
	// fmt.Println(body)
	json.Unmarshal(body, &response)
	// fmt.Println(response)
	for _, node := range response {
		if node.State == "leader" {
			host = strings.Split(node.ClientURL, "//")[1]
		}
	}
	return host
}

func publishService(client *etcd.Client, host string, ttl uint64) {
	for {
		setEtcd(client, swarmetcd, host, ttl)
		time.Sleep(timeout)
	}
}

func setEtcd(client *etcd.Client, key, value string, ttl uint64) {
	_, err := client.Set(key, value, ttl)
	if err != nil {
		log.Println(err)
	}
}
func main() {
	etcdproto := "etcd://" + getleaderHost() + swarmpath
	etcdhost := os.Getenv("HOST")
	addr := "--addr=" + etcdhost + ":2375"
	client := etcd.NewClient([]string{"http://" + etcdhost + ":" + etcdport})
	var wd, _ = os.Getwd()
	switch os.Args[1] {
	case "join":
		// fmt.Println("join")
		run(wd + "/deis-swarm join " + addr + " " + etcdproto)
	case "manage":
		// fmt.Println("manage")
		go publishService(client, etcdhost, uint64(ttl.Seconds()))
		run(wd + "/deis-swarm manage " + etcdproto)
	}
}
