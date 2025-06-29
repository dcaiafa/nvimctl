package main

import (
	"fmt"
	"os"
)

func main() {
	err := CmdRoot().Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
