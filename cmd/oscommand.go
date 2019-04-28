package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type (
	// OSCommand represents a command that runs on the OS
	OSCommand struct {
		Args        []string
		Cmd         string
		Description string
		SilentError bool
	}
)

// Run executes the oscommand
func (c *OSCommand) Run() {
	var b bytes.Buffer
	Cmd := exec.Command(c.Cmd, c.Args...)
	Cmd.Stdout = &b
	Cmd.Stderr = &b
	if err := Cmd.Run(); !c.SilentError && err != nil {
		fmt.Println(fmt.Errorf("%s: %v", c.Description, err))
		fmt.Println(b.String())
		os.Exit(1)
	} else {
		fmt.Println(c.Description)
		fmt.Println(b.String())
	}
}
