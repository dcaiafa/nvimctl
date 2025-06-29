package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func cmdCD() *cobra.Command {
	c := &cobra.Command{
		Use:  "cd path",
		Args: cobra.ExactArgs(1),
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

		return nv.ExecLua(
			"local target_path = (...); vim.api.nvim_set_current_dir(target_path)",
			nil, targetPath)
	}

	return c
}
