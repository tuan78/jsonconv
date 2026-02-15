package cli

import (
	"flag"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// version can be set at build time using ldflags:
	// go build -ldflags="-X 'github.com/tuan78/jsonconv/cmd.version=v1.2.3'"
	version = "dev"
)

type RootFlags struct {
	InputPath  string
	OutputPath string
	RawData    string
}

var rootFlags = &RootFlags{}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "jsonconv",
		Short:   "Tool for flattening JSON and converting JSON to CSV",
		Long:    "Tool for flattening JSON and converting JSON to CSV",
		Version: version,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			flag.Parse()
		},
	}
	// Add flags.
	cmd.PersistentFlags().StringVarP(&rootFlags.RawData, "data", "d", "", "raw JSON data. If both '--data' and '--in' are not set, reads from Stdin instead")
	cmd.PersistentFlags().StringVarP(&rootFlags.InputPath, "in", "i", "", "input file path. If both '--data' and '--in' are not set, reads from Stdin instead")
	cmd.PersistentFlags().StringVarP(&rootFlags.OutputPath, "out", "o", "", "output file path. It not set, prints to Stdout instead")

	// Add commands.
	cmd.AddCommand(NewFlattenCmd())
	cmd.AddCommand(NewCsvCmd())

	// Parse flags.
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	return cmd
}
