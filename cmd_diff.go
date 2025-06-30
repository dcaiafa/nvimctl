package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func cmdDiff() *cobra.Command {
	c := &cobra.Command{
		Use:   "diff <file1> <file2>",
		Short: "Compare two files in diff mode",
		Long:  `Opens two files in a new tab with diff mode enabled. Blocks until the tab is closed.`,
		Args:  cobra.ExactArgs(2),
	}

	c.RunE = func(cmd *cobra.Command, args []string) error {
		nv, err := DialNvim()
		if err != nil {
			return fmt.Errorf("failed to connect to nvim: %w", err)
		}

		// Get absolute paths for both files.
		absPath1, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("failed to get absolute path for file 1: %w", err)
		}
		absPath2, err := filepath.Abs(args[1])
		if err != nil {
			return fmt.Errorf("failed to get absolute path for file 2: %w", err)
		}

		// Create a channel to wait for tab close notification.
		done := make(chan bool)

		// Register handler for the notification that the operation is complete.
		nv.RegisterHandler("nvimctl_diff_done", func(args ...any) {
			done <- true
		})

		// Format the command to trigger the notification from the editor.
		notifyRPC := fmt.Sprintf(`call rpcnotify(%v, "nvimctl_diff_done")`, nv.ChannelID())

		// Escape the file paths to be used in commands.
		escaped, err := NvimEscapeSlice(nv, []string{absPath1, absPath2})
		if err != nil {
			return fmt.Errorf("failed to escape paths: %w", err)
		}
		escapedPath1 := escaped[0]
		escapedPath2 := escaped[1]

		// Get the current buffer number so we can restore it when we are done.
		var termBuf int
		err = nv.Eval(`bufnr("%")`, &termBuf)
		if err != nil {
			return fmt.Errorf("failed to get terminal buffer number: %w", err)
		}

		var (
			tabn int // Tab used for diff.
			bufl int // Buffer with the left diff file.
			bufr int // Buffer with the right diff file.
		)

		// Create the tab for diff'ing.
		batch := nv.NewBatch()
		batch.Command("tabnew")
		batch.Eval(`tabpagenr()`, &tabn)

		// Set up right diff.
		batch.Command(fmt.Sprintf("edit %s", escapedPath2))
		batch.Command("diffthis")
		batch.Command("setlocal bufhidden=wipe")
		batch.Command("augroup nvimctl_diff")
		batch.Command("au!")
		batch.Command("au BufDelete <buffer> " + notifyRPC)
		batch.Eval(`bufnr("%")`, &bufr)

		// Set up left diff.
		batch.Command(fmt.Sprintf("vert diffsplit %s", escapedPath1))
		batch.Command("au BufDelete <buffer> " + notifyRPC)
		batch.Command("setlocal bufhidden=wipe")
		batch.Eval(`bufnr("%")`, &bufl)
		batch.Command("augroup END")

		err = batch.Execute()
		if err != nil {
			return fmt.Errorf("failed to setup diff: %w", err)
		}

		// Wait for either file to be closed.
		<-done

		// Clean up.
		batch = nv.NewBatch()
		batch.Command("augroup nvimctl_diff")
		batch.Command("au!")
		batch.Command("augroup END")
		batch.Command(fmt.Sprintf("silent! bdel %v", bufl))
		batch.Command(fmt.Sprintf("silent! bdel %v", bufr))
		batch.Command(fmt.Sprintf("silent! tabclose %v", tabn))
		err = batch.Execute()
		if err != nil {
			return fmt.Errorf("failed to clean up: %w", err)
		}

		// Restore terminal buffer.
		batch = nv.NewBatch()
		batch.Command(fmt.Sprintf("buffer %d", termBuf))
		batch.Command("startinsert")
		err = batch.Execute()
		if err != nil {
			return fmt.Errorf("failed to restore terminal: %w", err)
		}

		return nil
	}

	return c
}
