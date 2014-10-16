package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fi, _ := os.Stdin.Stat()
	if fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("this app only works using the stdout of another process")
		os.Exit(1)
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var input map[string]interface{}
	err = json.Unmarshal(bytes, &input)

	if err != nil {
		fmt.Println("invalid application json configuration")
		os.Exit(1)
	}

	if input["release"] == nil {
		fmt.Println("invalid application version")
		os.Exit(1)
	}

	release := input["release"].(map[string]interface{})
	fmt.Println(release["version"])
}
