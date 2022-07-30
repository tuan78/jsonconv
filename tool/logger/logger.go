package logger

import (
	"fmt"

	"github.com/spf13/cobra"
)

type (
	// A Logger prints formatted string to desired outputs (stdout, stdin,
	// stderr or byte buffer) that can improve the testability.
	Logger interface {
		Printf(format string, i ...interface{})
	}

	logger struct {
		cmd *cobra.Command
	}
)

func NewLogger(cmd *cobra.Command) Logger {
	return &logger{cmd: cmd}
}

func (l *logger) Printf(format string, i ...interface{}) {
	fmt.Fprintf(l.cmd.OutOrStdout(), format, i...)
}
