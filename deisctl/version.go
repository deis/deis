package main

import (
    "fmt"

    "github.com/deis/deis/version"
)

var cmdVersion = &Command{
    Name:        "version",
    Description: "Print the version and exit",
    Summary:     "Print the version and exit",
    Run:         runVersion,
}

func runVersion(args []string) (exit int) {
    fmt.Println(version.Version)
    return
}
