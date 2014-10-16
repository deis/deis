package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/moraes/config"
)

func main() {
	fi, _ := os.Stdin.Stat()
	if fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("this app only works using the stdout of another process")
		os.Exit(1)
	}

	bytes, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Println("invalid input")
		os.Exit(1)
	}

	cfg, err := config.ParseYaml(string(bytes))

	if err != nil {
		fmt.Println("the procfile does not contains a valid yaml structure")
		os.Exit(1)
	}

	toSring, _ := json.Marshal(cfg.Root)
	fmt.Println(string(toSring))
}
