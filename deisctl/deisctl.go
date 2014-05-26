package main

import (
    "flag"
    "fmt"
    "os"
    "text/tabwriter"
)

const (
    cliName        = "deisctl"
    cliDescription = "deisctl is a command-line interface to Deis, the open-source Application Deployment Platform."
)

var (
    commands []*Command
    globalFlagset *flag.FlagSet = flag.NewFlagSet("deisctl", flag.ExitOnError)
    out *tabwriter.Writer

    // flags used by all commands
    globalFlags = struct {
        Version bool
    }{}
)

type Command struct {
    Name        string       // Name of the Command and the string to use to invoke it
    Summary     string       // One-sentence summary of what the Command does
    Usage       string       // Usage options/arguments
    Description string       // Detailed description of command
    Flags       flag.FlagSet // Set of flags associated with this command

    Run func(args []string) int // Run a command with the given arguments, return exit status
}

func init() {
    globalFlagset.BoolVar(&globalFlags.Version, "version", false, "Print the version and exit")

    out = new(tabwriter.Writer)
    out.Init(os.Stdout, 0, 8, 1, '\t', 0)
    commands = []*Command{
        cmdHelp,
        cmdVersion,
    }
}

func getAllFlags() (flags []*flag.Flag) {
    return getFlags(globalFlagset)
}

func getFlags(flagset *flag.FlagSet) (flags []*flag.Flag) {
    flags = make([]*flag.Flag, 0)
    flagset.VisitAll(func(f *flag.Flag) {
        flags = append(flags, f)
    })
    return
}

func main() {
    // parse global arguments
    globalFlagset.Parse(os.Args[1:])

    var args = globalFlagset.Args()

    // no command specified - trigger help
    if len(args) < 1 {
        args = append(args, "help")
    }

    // deal specially with --version
    if globalFlags.Version {
        args[0] = "version"
    }

    var cmd *Command

    // determine which Command should be run
    for _, c := range commands {
        if c.Name == args[0] {
            cmd = c
            if err := c.Flags.Parse(args[1:]); err != nil {
                fmt.Println(err.Error())
                os.Exit(2)
            }
            break
        }
    }

    if cmd == nil {
        fmt.Printf("%v: unknown subcommand: %q\n", cliName, args[0])
        fmt.Printf("Run '%v help' for usage.\n", cliName)
        os.Exit(2)
    }

    os.Exit(cmd.Run(cmd.Flags.Args()))
}
