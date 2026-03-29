package main

import (
	"fmt"
	"os"

	"github.com/tuan78/jsonconv/v2/internal/cli"
)

var exitFn = os.Exit

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Command execution failed, err: %v\n", err)
		exitFn(1)
	}
}
