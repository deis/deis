package types

import (
	"io"
	"os/exec"
)

// Command struct to execute commands.
type Command struct {
	Stdout, Stderr io.Writer

	cmd *exec.Cmd
}
