package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func cmdEdit() *cobra.Command {
	c := &cobra.Command{
		Use:   "edit <file>",
		Short: "Edit a file in a new split window",
		Long:  `Opens a file in a new split window and blocks until the window or buffer is closed.`,
		Args:  cobra.ExactArgs(1),
	}

	c.RunE = func(cmd *cobra.Command, args []string) error {
		nv, err := DialNvim()
		if err != nil {
			return fmt.Errorf("failed to connect to nvim: %w", err)
		}

		filePath := args[0]
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		// Leave terminal mode if necessary.
		err = nv.FeedKeys("\x1c\x0e", "n", false)
		if err != nil {
			return fmt.Errorf("failed to leave terminal mode: %w", err)
		}

		// Create a channel to wait for buffer close notification.
		done := make(chan bool)

		// Register handler for notifications.
		nv.RegisterHandler("nvimctl_done", func(args ...any) {
			done <- true
		})

		// Escape the file path.
		escaped, err := NvimEscape(nv, []string{absPath})
		if err != nil {
			return fmt.Errorf("failed to escape path: %w", err)
		}
		escapedPath := escaped[0]

		batch := nv.NewBatch()

		// Open the file in a new split window.
		batch.Command(fmt.Sprintf("below new %s", escapedPath))

		// Set up auto commands for the new buffer.
		batch.Command("augroup nvimctl_edit")
		batch.Command(`au!`)

		// Notify if the user closes the buffer.
		batch.Command(fmt.Sprintf(
			`au BufDelete <buffer> call rpcnotify(%v, "nvimctl_done")`, nv.ChannelID()))

		// If the user writes to the file, close the buffer and notify.
		batch.Command(fmt.Sprintf(
			`au BufWritePost <buffer> bw | call rpcnotify(%v, "nvimctl_done")`,
			nv.ChannelID()))
		batch.Command("augroup END")

		err = batch.Execute()
		if err != nil {
			return fmt.Errorf("failed to setup buffer: %w", err)
		}

		// Wait for the buffer to close.
		<-done

		// Enter insert mode after getting back to the terminal.
		nv.Command("startinsert")

		return nil
	}

	return c
}
