# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

nvimctl is a Go-based CLI tool for controlling Neovim from the command line. It connects to a parent Neovim instance via the NVIM environment variable and executes commands through Neovim's Lua API.

## Development Commands

Standard Go commands are used for development:
- `go build` - Build the nvimctl binary
- `go run . [command]` - Run without building
- `go mod tidy` - Clean up module dependencies
- `go test ./...` - Run tests (no tests currently exist) - **Use this for verification instead of `go build` to avoid creating executables**
- `go fmt ./...` - Format code

## Architecture

### Command Structure
The project uses the Cobra library for CLI commands. Each command:
1. Lives in its own file (`cmd_<name>.go`)
2. Has a function `cmd<Name>()` that returns `*cobra.Command`
3. Is registered in `cmd_root.go` via `c.AddCommand(cmd<Name>())`

### Adding New Commands
To add a new command:
1. Create `cmd_yourcommand.go`. Use `cmd_cd.go` as a reference.
2. Register in `cmd_root.go`: `c.AddCommand(cmdYourCommand())`

### Neovim Connection and Command Execution
- Requires the `NVIM` environment variable (set automatically when running from Neovim terminal)
- Connection established via `DialNvim()` in `nvim.go`
- Use `NvimEscape` and `NvimEscapeSlice` to escape nvim arguments when necessary.
- Use `cmd_diff.go` as a reference for commands that block.

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
