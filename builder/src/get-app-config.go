package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/deis/deis/builder"
)

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
		os.Exit(1)
	}

	if *builderKey == "" {
		fmt.Println("invalid builder key")
		os.Exit(1)
	}

	if *user == "" {
		fmt.Println("invalid user")
		os.Exit(1)
	}

	if *app == "" {
		fmt.Println("invalid app")
		os.Exit(1)
	}

	data, err := json.Marshal(&builder.ConfigHook{ReceiveUser: *user, ReceiveRepo: *app})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b := bytes.NewReader(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", *url, b)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-Deis-Builder-Auth", *builderKey)

	res, err := client.Do(req)
	defer res.Body.Close()

	if res.StatusCode == 404 {
		fmt.Println("Check the Controller. Is it running?")
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err != nil || res.StatusCode != 200 {
		fmt.Println("failed retrieving config from controller")
		fmt.Printf("%v\n", body)
		os.Exit(1)
	}

	config, err := builder.ParseConfig(body)
	if err != nil {
		fmt.Println("failed parsing config from controller")
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	toString, err := json.Marshal(config)
	fmt.Println(string(toString))
}
