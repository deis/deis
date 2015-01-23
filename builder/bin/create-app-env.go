package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deis/deis/builder"
)

func main() {
	directory := flag.String("directory", "", "Application environment directory")

	flag.Parse()

	if flag.NFlag() < 1 {
		fmt.Printf("Usage: -directory\n")
		os.Exit(1)
	}

	if fi, _ := os.Stdin.Stat(); fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("this app only works using the stdout of another process")
		os.Exit(1)
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err = builder.CreateEnvFilesFomConfig(*directory, bytes); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
