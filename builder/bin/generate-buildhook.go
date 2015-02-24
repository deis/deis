package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/deis/deis/builder"
)

func assert(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func usage(s string) {
	fmt.Printf("Usage: %s <sha> <receive_user> <receive_repo> <image> <procfile> <dockerfile>\n", s)
}

func main() {
	if len(os.Args) != 7 {
		usage(os.Args[0])
		os.Exit(1)
	}

	var procfile builder.ProcessType
	assert(json.Unmarshal([]byte(os.Args[5]), &procfile))

	var dockerfile = os.Args[6]
	if dockerfile == "false" {
		dockerfile = ""
	}

	buildHook := builder.BuildHook{
		Sha:         os.Args[1],
		ReceiveUser: os.Args[2],
		ReceiveRepo: os.Args[3],
		Image:       os.Args[4],
		Procfile:    procfile,
		Dockerfile:  dockerfile,
	}

	b, err := json.Marshal(buildHook)
	assert(err)

	fmt.Println(string(b))
}
