package updatectl

import (
	"flag"
	"fmt"
	"github.com/coreos/updatectl/auth"
	"github.com/coreos/updatectl/client/update/v1"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
)

const (
	OK = iota
	// Error Codes
	ERROR_API
	ERROR_USAGE
	ERROR_NO_COMMAND

	cliName        = "deisctl"
	cliDescription = "deisctl update is a command line driven interface to the roller."
)

type StringFlag struct {
	value    *string
	required bool
}

func (f *StringFlag) Set(value string) error {
	f.value = &value
	return nil
}

func (f *StringFlag) Get() *string {
	return f.value
}

func (f *StringFlag) String() string {
	if f.value != nil {
		return *f.value
	}
	return ""
}

type Command struct {
	Name        string       // Name of the Command and the string to use to invoke it
	Summary     string       // One-sentence summary of what the Command does
	Usage       string       // Usage options/arguments
	Description string       // Detailed description of command
	Flags       flag.FlagSet // Set of flags associated with this command
	Run         handlerFunc  // Run a command with the given arguments
	Subcommands []*Command   // Subcommands for this command.
}

var (
	out           *tabwriter.Writer
	globalFlagSet *flag.FlagSet
	commands      []*Command
	globalFlags   struct {
		Server        string
		User          string
		Key           string
		Debug         bool
		Version       bool
		Help          bool
		SkipSSLVerify bool
	}
)

func init() {

	out = new(tabwriter.Writer)
	out.Init(os.Stdout, 0, 8, 1, '\t', 0)
	server := "http://localhost:8000" // default server
	if serverEnv := os.Getenv("UPDATECTL_SERVER"); serverEnv != "" {
		server = serverEnv
	}

	globalFlagSet = flag.NewFlagSet(cliName, flag.ExitOnError)
	globalFlagSet.StringVar(&globalFlags.Server, "server", server, "Update server to connect to")
	globalFlagSet.BoolVar(&globalFlags.Debug, "debug", false, "Output debugging info to stderr")
	globalFlagSet.BoolVar(&globalFlags.Version, "version", false, "Print version information and exit.")
	globalFlagSet.BoolVar(&globalFlags.Help, "help", false, "Print usage information and exit.")
	globalFlagSet.BoolVar(&globalFlags.SkipSSLVerify, "skip-ssl-verify", false, "Don't check SSL certificates.")
	globalFlagSet.StringVar(&globalFlags.User, "user", os.Getenv("UPDATECTL_USER"), "API Username")
	globalFlagSet.StringVar(&globalFlags.Key, "key", os.Getenv("UPDATECTL_KEY"), "API Key")

	commands = []*Command{
		cmdInstance,
	}
}

type handlerFunc func([]string, *update.Service, *tabwriter.Writer) int

func getHawkClient(user string, key string) *http.Client {
	return &http.Client{
		Transport: &auth.HawkRoundTripper{
			User:          user,
			Token:         key,
			SkipSSLVerify: globalFlags.SkipSSLVerify,
		},
	}
}

func handle(fn handlerFunc) func(f *flag.FlagSet) int {
	return func(f *flag.FlagSet) (exit int) {
		user := globalFlags.User
		key := globalFlags.Key
		client := getHawkClient(user, key)
		service, err := update.New(client)
		if err != nil {
			log.Fatal(err)
		}
		service.BasePath = globalFlags.Server + "/_ah/api/update/v1/"
		exit = fn(f.Args(), service, out)
		return
	}
}

func getAllFlags() (flags []*flag.Flag) {
	return getFlags(globalFlagSet)
}

func getFlags(flagset *flag.FlagSet) (flags []*flag.Flag) {
	flags = make([]*flag.Flag, 0)
	flagset.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f)
	})
	return
}

// determine which Command should be run
func findCommand(search string, args []string, commands []*Command) (cmd *Command, name string) {
	if len(args) < 1 {
		return
	}
	if search == "" {
		search = args[0]
	} else {
		search = fmt.Sprintf("%s %s", search, args[0])
	}
	name = search
	for _, c := range commands {
		if c.Name == search {
			cmd = c
			if errHelp := c.Flags.Parse(args[1:]); errHelp != nil {
				//printCommandUsage(cmd)
				os.Exit(ERROR_USAGE)
			}
			if len(cmd.Subcommands) != 0 {
				subArgs := cmd.Flags.Args()
				var subCmd *Command
				subCmd, name = findCommand(search, subArgs, cmd.Subcommands)
				if subCmd != nil {
					cmd = subCmd
				}
			}
			break
		}
	}
	return
}

func Update(Args []string) {
	globalFlagSet.Parse(Args)
	var args = globalFlagSet.Args()

	if len(args) < 1 {
		args = append(args, "help")
	}

	cmd, name := findCommand("", args, commands)

	if cmd == nil {
		fmt.Printf("%v: unknown subcommand: %q\n", cliName, name)
		fmt.Printf("Run '%v help' for usage.\n", cliName)
		os.Exit(ERROR_NO_COMMAND)
	}
	if cmd.Run == nil {
		//	printCommandUsage(cmd)
		os.Exit(ERROR_USAGE)
	} else {
		exit := handle(cmd.Run)(&cmd.Flags)
		if exit == ERROR_USAGE {
			//	printCommandUsage(cmd)
		}
		os.Exit(exit)
	}
}
