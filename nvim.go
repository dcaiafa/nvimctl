package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/neovim/go-client/nvim"
)

func DialNvim() (*nvim.Nvim, error) {
	nvimAddr := os.Getenv("NVIM")
	if nvimAddr == "" {
		return nil, fmt.Errorf("could not detect parent nvim")
	}

	return nvim.Dial(nvimAddr)
}

type Command struct {
	Command string
	Args    []string
	Output  bool
}

func NvimExec(nv *nvim.Nvim, cmd *Command) (string, error) {
	batch := nv.NewBatch()

	cmdParts := make([]string, 0, len(cmd.Args)+1)
	cmdParts = append(cmdParts, cmd.Command)

	for _, arg := range cmd.Args {
		cmdParts = append(cmdParts, "")
		batch.Call("fnameescape", &cmdParts[len(cmdParts)-1], arg)
	}
	err := batch.Execute()
	if err != nil {
		return "", fmt.Errorf("failed to escape arguments: %w", err)
	}

	fullCmd := strings.Join(cmdParts, " ")

	output, err := nv.Exec(fullCmd, cmd.Output)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return output, nil
}
