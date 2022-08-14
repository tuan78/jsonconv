package main

import (
	"fmt"
	"os"

	"github.com/tuan78/jsonconv/cmd"
)

func main() {
	// Execute command.
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Command execution failed, err: %v\n", err)
		os.Exit(1)
	}
}
