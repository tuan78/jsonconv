package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/tuan78/jsonconv/tool/cmd"
	"github.com/tuan78/jsonconv/tool/params"
)

var (
	version = "v0.1.0"
)

var (
	rootCmd = &cobra.Command{
		Use:     "jsonconv",
		Short:   "Tool for flattening and converting JSON.",
		Long:    "Tool for flattening and converting JSON (JSON to CSV, JSON from CSV, JSON from Excel, and more).",
		Version: version,
	}
)

func main() {
	// Add flags.
	rootCmd.PersistentFlags().StringVarP(&params.RawData, "data", "d", "", "raw JSON array data")
	rootCmd.PersistentFlags().StringVarP(&params.InputPath, "in", "i", "", "input file path")
	rootCmd.PersistentFlags().StringVarP(&params.OutputPath, "out", "o", "", "output file path")

	// Add commands.
	rootCmd.AddCommand(cmd.NewFlattenCmd())
	rootCmd.AddCommand(cmd.NewCsvCmd())

	// Parse flags.
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()

	// Execute command.
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Command execution failed, error: %v", err))
		os.Exit(1)
	}
}
