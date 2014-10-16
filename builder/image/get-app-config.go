package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/franela/goreq"
)

// Config repo from which extract the configuration and user to use
type Config struct {
	ReceiveUser string `json:"receive_user"`
	ReceiveRepo string `json:"receive_repo"`
}

// Release application configuration
type Release struct {
	Owner   string                 `json:"owner"`
	App     string                 `json:"app"`
	Values  map[string]interface{} `json:"values"`
	Memory  map[string]interface{} `json:"memory"`
	CPU     map[string]interface{} `json:"cpu"`
	Tags    map[string]interface{} `json:"tags"`
	UUID    string                 `json:"uuid"`
	Created time.Time              `json:"created"`
	Updated time.Time              `json:"updated"`
}

const (
	contentType string = "application/json"
	userAgent   string = "deis-builder"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: [options]\n\n")
		flag.PrintDefaults()
	}
}

func main() {
	url := flag.String("url", "", "Controller hook URL")
	builderKey := flag.String("key", "", "Builder Key")
	user := flag.String("user", "", "Controller username")
	app := flag.String("app", "", "Controller application name")

	flag.Parse()

	if flag.NFlag() < 4 {
		flag.Usage()
		os.Exit(1)
	}

	if *url == "" {
		fmt.Println("invalid url")
		os.Exit(64)
	}

	if *builderKey == "" {
		fmt.Println("invalid builder key")
		os.Exit(64)
	}

	if *user == "" {
		fmt.Println("invalid user")
		os.Exit(64)
	}

	if *app == "" {
		fmt.Println("invalid app")
		os.Exit(64)
	}

	data := Config{ReceiveUser: *user, ReceiveRepo: *app}

	req := goreq.Request{
		Method:      "POST",
		Uri:         *url,
		Body:        data,
		ContentType: contentType,
		Accept:      contentType,
		UserAgent:   userAgent,
	}

	req.AddHeader("X-Deis-Builder-Auth", *builderKey)

	res, err := req.Do()

	if res.StatusCode == 404 {
		fmt.Println("check the deis-controller. Is not running")
		os.Exit(2)
	}

	if err != nil || res.StatusCode != 200 {
		fmt.Println("failed retrieving config from controller")
		body, _ := res.Body.ToString()
		fmt.Println(body)
		os.Exit(1)
	}

	var release Release
	res.Body.FromJsonTo(&release)
	toSring, _ := json.Marshal(release)
	fmt.Println(string(toSring))
}
