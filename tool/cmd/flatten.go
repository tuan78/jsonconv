package cmd

import "github.com/spf13/cobra"

func NewFlattenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flatten",
		Short: "Flatten JSON object and JSON array",
		Long:  "Flatten JSON object and JSON array",
		RunE: func(cmd *cobra.Command, args []string) error {
			return processFlattenCmd()
		},
	}
	return cmd
}

func processFlattenCmd() error {
	return nil
}
