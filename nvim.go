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

func NvimEscape(nv *nvim.Nvim, vs []string) ([]string, error) {
	batch := nv.NewBatch()

	evs := make([]string, len(vs))
	for i, v := range vs {
		batch.Call("fnameescape", &evs[i], v)
	}
	err := batch.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to escape values: %w", err)
	}

	return evs, nil
}

func NvimExec(nv *nvim.Nvim, cmd *Command) (string, error) {
	cmdParts := append([]string{cmd.Command}, cmd.Args...)
	cmdParts, err := NvimEscape(nv, cmdParts)
	if err != nil {
		return "", err
	}

	fullCmd := strings.Join(cmdParts, " ")
	output, err := nv.Exec(fullCmd, cmd.Output)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return output, nil
}

func LeaveTerminal(nv *nvim.Nvim) error {
	// Leave terminal mode if necessary (equivalent to feedkeys Ctrl-\ Ctrl-N).
	err := nv.FeedKeys("\x1c\x0e", "n", false)
	if err != nil {
		return fmt.Errorf("failed to leave terminal mode: %w", err)
	}
	return nil
}
