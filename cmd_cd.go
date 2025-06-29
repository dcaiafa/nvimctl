package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func cmdCD() *cobra.Command {
	c := &cobra.Command{
		Use:   "cd <path>",
		Short: "Change neovim's current directory",
		Args:  cobra.ExactArgs(1),
	}

	c.RunE = func(cmd *cobra.Command, args []string) error {
		targetPath := args[0]

		targetPath, err := filepath.Abs(targetPath)
		if err != nil {
			return fmt.Errorf("path is invalid")
		}

		nv, err := DialNvim()
		if err != nil {
			return nil
		}

		_, err = NvimExec(nv, &Command{
			Command: "cd",
			Args:    []string{targetPath},
		})

		return err
	}

	return c
}
