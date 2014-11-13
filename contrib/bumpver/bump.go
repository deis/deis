package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	docopt "github.com/docopt/docopt-go"
)

// Command processes command-line arguments according to a usage string specification.
func Command(argv []string) (returnCode int) {
	usage := `
Updates the current semantic version number to a new one in the specified
source and doc files.

Usage:
  bumpver [-f <current>] <version> <files>...

Options:
  -f --from=<current>  An explicit version string to replace. Otherwise, use
                       the first semantic version found in the first file.
 `
	argv = parseArgs(argv)
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "bumpversion 0.1.0", true, false)
	if err != nil {
		return 1
	}
	if len(args) == 0 {
		return 0
	}
	re := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3})`)
	// validate that <version> is a proper semver string
	version := (args["<version>"].(string))
	if version != re.FindString(version) {
		return onError(fmt.Errorf("Error: '%s' is not a valid semantic version string", version))
	}
	files := args["<files>"].([]string)
	var from []byte
	var data []byte
	if args["--from"] != nil {
		from = []byte(args["--from"].(string))
	}
	for _, name := range files {
		src, err := ioutil.ReadFile(name)
		if err != nil {
			return onError(err)
		}
		if len(from) == 0 {
			// find the first semver match in the file, if any
			from = re.Find(src)
			if from = re.Find(src); len(from) == 0 {
				fmt.Printf("Skipped %s\n", name)
				continue
			}
		}
		// replace all occurrences in source and doc files
		data = bytes.Replace(src, from, []byte(version), -1)
		f, err := os.Create(name)
		if err != nil {
			return onError(err)
		}
		if _, err = f.Write(data); err != nil {
			return onError(err)
		}
		fmt.Printf("Bumped %s\n", name)
	}
	return 0
}

func onError(err error) int {
	fmt.Println(err.Error())
	return 1
}

func parseArgs(argv []string) []string {
	if argv == nil {
		argv = os.Args[1:]
	}

	if len(argv) > 0 {
		// parse "help <command>" as "command> --help"
		if argv[0] == "help" {
			argv = append(argv[1:], "--help")
		}
	}
	return argv
}
