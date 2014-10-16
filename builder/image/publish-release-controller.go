package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/franela/goreq"
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
		os.Exit(64)
	}

	if *builderKey == "" {
		fmt.Println("invalid builder key")
		os.Exit(64)
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("invalid json payload")
		os.Exit(64)
	}

	postBody := strings.Replace(string(bytes), "'", "", -1)

	// Check for a variable trying to exploit Shellshock.
	potencialExploit := regexp.MustCompile(`\(\)\s+\{[^\}]+\};\s+(.*)`)
	if potencialExploit.MatchString(postBody) {
		fmt.Println("")
		fmt.Println("ATTENTION: an environment variable in the app is trying to exploit Shellshock. Aborting...")
		fmt.Println("")
		os.Exit(1)
	}

	req := goreq.Request{
		Method:      "POST",
		Uri:         *url,
		Body:        postBody,
		ContentType: contentType,
		Accept:      contentType,
		UserAgent:   userAgent,
	}

	req.AddHeader("X-Deis-Builder-Auth", *builderKey)

	res, err := req.Do()

	// Read json response from body
	body, _ := res.Body.ToString()
	var response map[string]interface{}
	jsonErr := json.Unmarshal([]byte(body), &response)

	if jsonErr != nil {
		fmt.Println("invalid controller json response")
		fmt.Println(body)
		os.Exit(1)
	}

	if err != nil || res.StatusCode != 200 {
		fmt.Println("failed retrieving config from controller")
		fmt.Println(response["detail"].(string))
		os.Exit(1)
	}

	if res.StatusCode == 404 {
		fmt.Println("check the deis-controller. Is not running")
		os.Exit(2)
	}

	toSring, _ := json.Marshal(response)
	fmt.Println(string(toSring))
	os.Exit(0)
}
