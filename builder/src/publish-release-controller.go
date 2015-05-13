package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

	flag.Parse()

	if flag.NFlag() < 2 {
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

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("invalid json payload")
		os.Exit(1)
	}

	postBody := strings.Replace(string(bytes), "'", "", -1)

	// Check for a variable trying to exploit Shellshock.
	potentialExploit := regexp.MustCompile(`\(\)\s+\{[^\}]+\};\s+(.*)`)
	if potentialExploit.MatchString(postBody) {
		fmt.Println("")
		fmt.Println("ATTENTION: an environment variable in the app is trying to exploit Shellshock. Aborting...")
		fmt.Println("")
		os.Exit(1)
	}

	b := strings.NewReader(postBody)
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

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if res.StatusCode == 503 {
		log.Fatalln("check the controller. is it running?")
	} else if res.StatusCode != 200 {
		log.Fatalf("failed retrieving config from controller: %s\n", body)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("invalid controller json response")
		fmt.Println(string(body))
		os.Exit(1)
	}

	toString, _ := json.Marshal(response)
	fmt.Println(string(toString))
	os.Exit(0)
}
