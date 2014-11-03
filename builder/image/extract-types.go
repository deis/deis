package main

import (
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
	// this should not happen but just in case return an empty type
	if err != nil {
		fmt.Println("{}")
		os.Exit(0)
	}

	cfg, err := config.ParseYaml(string(bytes))

	if err != nil {
		fmt.Println("the procfile does not contains a valid yaml structure")
		os.Exit(1)
	}

	defaultType, err := cfg.Get("default_process_types")

	// some buildpacks don't supply a default process type
	// as Heroku does not make them mandatory
	if err != nil {
		fmt.Println("{}")
		os.Exit(0)
	}

	fmt.Println(defaultType)
}
