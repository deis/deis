package types

import (
	"io"
	"os/exec"
)

type Command struct {
	Stdout, Stderr io.Writer

	cmd *exec.Cmd
}
