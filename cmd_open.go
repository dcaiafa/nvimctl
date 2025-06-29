package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func cmdOpen() *cobra.Command {
	c := &cobra.Command{
		Use:   "open <file> [wincmd]",
		Short: "Open a file in Neovim with optional window movement",
		Long: `Open a file in Neovim. Optionally specify window movements before opening.
		
Example window movements:
  h - move to left window
  j - move to bottom window
  k - move to top window
  l - move to right window
  
You can combine movements, e.g., "hjl" to move left, down, then right.`,
		Args: cobra.RangeArgs(1, 2),
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

		// Leave terminal mode if necessary (equivalent to feedkeys Ctrl-\ Ctrl-N).
		err = nv.FeedKeys("\x1c\x0e", "n", false)
		if err != nil {
			return fmt.Errorf("failed to leave terminal mode: %w", err)
		}

		// Apply window movements if provided.
		if len(args) > 1 {
			wincmd := args[1]
			for _, pos := range wincmd {
				_, err = NvimExec(nv, &Command{
					Command: "wincmd",
					Args:    []string{string(pos)},
				})
				if err != nil {
					return fmt.Errorf("failed to execute wincmd %c: %w", pos, err)
				}
			}
		}

		// Open the file using drop command.
		_, err = NvimExec(nv, &Command{
			Command: "drop",
			Args:    []string{absPath},
		})
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		return nil
	}

	return c
}