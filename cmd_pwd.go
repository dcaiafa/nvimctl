package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cmdPwd() *cobra.Command {
	c := &cobra.Command{
		Use:   "pwd",
		Short: "Print Neovim's current working directory",
		Args:  cobra.NoArgs,
	}

	c.RunE = func(cmd *cobra.Command, args []string) error {
		nv, err := DialNvim()
		if err != nil {
			return err
		}

		var cwd string
		err = nv.Eval("getcwd()", &cwd)
		if err != nil {
			return err
		}

		fmt.Println(cwd)
		return nil
	}

	return c
}