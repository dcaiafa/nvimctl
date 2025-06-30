package main

import "github.com/spf13/cobra"

func CmdRoot() *cobra.Command {
	c := &cobra.Command{
		Short: "neovimctl is a cli for neovim",
	}

	c.AddCommand(cmdCD())
	c.AddCommand(cmdOpen())
	c.AddCommand(cmdEdit())
	c.AddCommand(cmdDiff())
	c.AddCommand(cmdPwd())

	return c
}
