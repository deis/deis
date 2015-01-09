package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/docker/docker/registry"
)

const (
	etcdPath = "/deis/registry/auth"
)

func main() {
	var cfgFile = flag.String("dockercfg", "/home/core/.dockercfg", "optional path to the .dockercfg directory")
	flag.Parse()

	fmt.Printf("using dockercfg file %v\n", *cfgFile)

	etcdHost := getopt("ETCD_HOST", "127.0.0.1")
	etcdClient := etcd.NewClient([]string{"http://" + etcdHost + ":4001"})
	// create required structure in etcd
	etcdClient.CreateDir(etcdPath, 0)

	cfg, err := registry.LoadConfig(*cfgFile)
	if err != nil {
		fmt.Errorf("%v", err.Error())
	}

	for registryHost := range cfg.Configs {
		config, ok := cfg.Configs[registryHost]
		if !ok {
			continue
		}
		etcdBase := etcdPath

		// replace the docker registry url
		if strings.EqualFold("https://index.docker.io/v1/", registryHost) {
			registryHost = "registry.hub.docker.com"
		}

		fmt.Printf("adding authentication for docker registry %v\n", registryHost)
		etcdBase = etcdBase + "/" + registryHost
		etcdClient.Set(etcdBase+"/username", config.Username, 0)
		etcdClient.Set(etcdBase+"/password", config.Password, 0)
		etcdClient.Set(etcdBase+"/email", config.Email, 0)
		fmt.Print("done")
	}
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
