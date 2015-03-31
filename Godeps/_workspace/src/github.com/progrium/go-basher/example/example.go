package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/dustin/go-jsonpointer"
	"github.com/progrium/go-basher"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func jsonPointer(args []string) int {
	if len(args) == 0 {
		return 3
	}
	bytes, err := ioutil.ReadAll(os.Stdin)
	assert(err)
	var o map[string]interface{}
	assert(json.Unmarshal(bytes, &o))
	println(jsonpointer.Get(o, args[0]).(string))
	return 0
}

func reverse(args []string) int {
	bytes, err := ioutil.ReadAll(os.Stdin)
	assert(err)
	runes := []rune(strings.Trim(string(bytes), "\n"))
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	println(string(runes))
	return 0
}

func main() {
	bash, _ := basher.NewContext("/bin/bash", false)
	bash.ExportFunc("json-pointer", jsonPointer)
	bash.ExportFunc("reverse", reverse)
	bash.HandleFuncs(os.Args)

	bash.Source("bash/example.bash", Asset)
	status, err := bash.Run("main", os.Args[1:])
	assert(err)
	os.Exit(status)
}
