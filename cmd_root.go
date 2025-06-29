package main

import "github.com/spf13/cobra"

func CmdRoot() *cobra.Command {
	c := &cobra.Command{
		Short:   "Control cli for neovim",
		Long:    "",
		Example: "",
	}

	c.AddCommand(cmdCD())
	c.AddCommand(cmdOpen())
	c.AddCommand(cmdEdit())
	c.AddCommand(cmdDiff())

	return c
}
