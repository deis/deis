package parser

import (
	"fmt"
	"strconv"
)

// docopt expects commands to be in the proper format, but we split them apart for
// routing purposes, so the commands need to be recombined.
func combineCommand(argv []string) []string {
	if len(argv) > 1 {
		return append([]string{argv[0] + ":" + argv[1]}, argv[2:]...)
	}

	return nil
}

func safeGetValue(args map[string]interface{}, key string) string {
	if args[key] == nil {
		return ""
	}
	return args[key].(string)
}

func responseLimit(limit string) (int, error) {
	if limit == "" {
		return -1, nil
	}

	return strconv.Atoi(limit)
}

// PrintUsage runs if no matching command is found.
func PrintUsage() {
	fmt.Println("Found no matching command, try 'deis help'")
	fmt.Println("Usage: deis <command> [<args>...]")
}
