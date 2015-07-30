package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deis/deis/builder"
)

func main() {
	if fi, _ := os.Stdin.Stat(); fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("this app requires the stdout of another process")
		os.Exit(1)
	}

	bytes, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Println("invalid input")
		os.Exit(1)
	}

	procfile, err := builder.YamlToJSON(bytes)

	if err != nil {
		fmt.Println("the Procfile is not valid yaml")
		os.Exit(1)
	}

	fmt.Println(procfile)
}
