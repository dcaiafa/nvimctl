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
