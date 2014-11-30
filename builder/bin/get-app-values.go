package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deis/deis/builder"
)

func main() {
	if fi, _ := os.Stdin.Stat(); fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("this app only works using the stdout of another process")
		os.Exit(1)
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	values, err := builder.ParseControllerConfig(bytes)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, value := range values {
		fmt.Println(value)
	}
}
