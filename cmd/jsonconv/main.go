package main

import (
	"fmt"
	"os"

	"github.com/tuan78/jsonconv/cmd"
)

func main() {
	// Execute command.
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Command execution failed, err: %v", err))
		os.Exit(1)
	}
}
