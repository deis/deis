package pkg

// TrimToDashes takes a slice of strings (e.g. a
// command line) and returns everything after the first
// double dash (--), if any are present
func TrimToDashes(args []string) []string {
	for i, arg := range args {
		if arg == "--" {
			return args[i+1:]
		}
	}
	return args
}
