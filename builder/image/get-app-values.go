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

	var release map[string]interface{}
	err = json.Unmarshal(bytes, &release)

	if err != nil {
		fmt.Println("invalid application json configuration")
		os.Exit(1)
	}

	if release["values"] == nil {
		fmt.Println()
		os.Exit(0)
	}

	values := release["values"].(map[string]interface{})

	for k, v := range values {
		fmt.Print(" -e " + k + "=\"" + v.(string) + "\"")
	}
}
