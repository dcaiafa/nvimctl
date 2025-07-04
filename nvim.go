package main

import (
	"fmt"
	"os"

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

func NvimEscape(nv *nvim.Nvim, v string) (string, error) {
	evs, err := NvimEscapeSlice(nv, []string{v})
	if err != nil {
		return "", err
	}
	return evs[0], nil
}

func NvimEscapeSlice(nv *nvim.Nvim, vs []string) ([]string, error) {
	batch := nv.NewBatch()

	evs := make([]string, len(vs))
	for i, v := range vs {
		batch.Call("fnameescape", &evs[i], v)
	}
	err := batch.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to escape value: %w", err)
	}

	return evs, nil
}

func LeaveTerminal(nv *nvim.Nvim) error {
	// Leave terminal mode if necessary (equivalent to feedkeys Ctrl-\ Ctrl-N).
	err := nv.FeedKeys("\x1c\x0e", "n", false)
	if err != nil {
		return fmt.Errorf("failed to leave terminal mode: %w", err)
	}
	return nil
}
