# nvimctl

A command-line interface for controlling Neovim from the terminal.

## Overview

nvimctl is a Go-based CLI tool that allows you to control a parent Neovim
instance from the command line. It connects to Neovim via the `NVIM` environment
variable and executes commands through Neovim's Lua API.

## Installation

Download the [latest
release](https://github.com/dcaiafa/nvimctl/releases/latest), unpack, and copy
to your PATH.

Or install the latest using Go:

```bash
go install github.com/dcaiafa/nvimctl@latest
```

Or build from source:

```bash
git clone https://github.com/dcaiafa/nvimctl
cd nvimctl
go build
```

## Usage

nvimctl must be run from within a Neovim terminal session. The tool
automatically detects the parent Neovim instance through the `NVIM` environment
variable.

### Available Commands

#### `open <file> [wincmd]`
Open a file in Neovim.

```bash
# Open file in current window
nvimctl open main.go

# Open file in top-left window (h=left, k=up)
nvimctl open main.go hk
```
#### `edit <file>`
Edit a file in a new split window, blocking until it is done. The window closes
automatically when you save the file or close the buffer. 

This is intended to be used with the `EDITOR` environment variable, and similar
mechanisms.

```bash
nvimctl edit config.yaml
```

#### `diff <file1> <file2>`
Compare two files in diff mode. 

```bash
nvimctl diff old_version.go new_version.go
```

#### `cd <path>`
Change Neovim's current directory.

```bash
nvimctl cd ~/projects/myproject
```

## Requirements

- Go 1.19 or later (for building)
- Neovim 0.7.0 or later
- Must be run from within a Neovim terminal

## License

MIT
