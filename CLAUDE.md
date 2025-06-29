# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

nvimctl is a Go-based CLI tool for controlling Neovim from the command line. It connects to a parent Neovim instance via the NVIM environment variable and executes commands through Neovim's Lua API.

## Development Commands

Standard Go commands are used for development:
- `go build` - Build the nvimctl binary
- `go run . [command]` - Run without building
- `go mod tidy` - Clean up module dependencies
- `go test ./...` - Run tests (no tests currently exist)
- `go fmt ./...` - Format code

## Architecture

### Command Structure
The project uses the Cobra library for CLI commands. Each command:
1. Lives in its own file (`cmd_<name>.go`)
2. Has a function `cmd<Name>()` that returns `*cobra.Command`
3. Is registered in `cmd_root.go` via `c.AddCommand(cmd<Name>())`

### Adding New Commands
To add a new command:
1. Create `cmd_yourcommand.go`:
```go
func cmdYourCommand() *cobra.Command {
    c := &cobra.Command{
        Use:   "yourcommand [args]",
        Short: "Brief description",
        Args:  cobra.ExactArgs(1),
    }
    
    c.RunE = func(cmd *cobra.Command, args []string) error {
        nv, err := DialNvim()
        if err != nil {
            return fmt.Errorf("failed to connect to nvim: %w", err)
        }
        
        // Execute Neovim commands using NvimExec.
        _, err = NvimExec(nv, &Command{
            Command: "yourcommand",
            Args:    []string{arg1, arg2},
            Output:  false, // Set to true if you need command output.
        })
        
        return err
    }
    
    return c
}
```
2. Register in `cmd_root.go`: `c.AddCommand(cmdYourCommand())`

### Neovim Connection and Command Execution
- Requires the `NVIM` environment variable (set automatically when running from Neovim terminal)
- Connection established via `DialNvim()` in `nvim.go`
- Commands should use `NvimExec()` for executing Neovim commands with automatic argument escaping
- `NvimExec()` uses `fnameescape()` to properly escape arguments before execution

### Key Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/neovim/go-client` - Neovim RPC client

## Important Notes
- This tool must be run from within a Neovim terminal session
- Always return errors up the chain (use `RunE`, not `Run`)
- Use `NvimExec()` for all Neovim command execution to ensure proper argument escaping
- Follow the established naming pattern: `cmd_<name>.go` for command files

## Code Style
- All Go comments, including single-line comments, should end with a period.