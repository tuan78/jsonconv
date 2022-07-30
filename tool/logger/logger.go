package logger

import (
	"fmt"

	"github.com/spf13/cobra"
)

type (
	CmdLogger interface {
		Printf(format string, i ...interface{})
	}

	cmdLogger struct {
		cmd *cobra.Command
	}
)

func NewCmdLogger(cmd *cobra.Command) CmdLogger {
	return &cmdLogger{cmd: cmd}
}

func (l *cmdLogger) Printf(format string, i ...interface{}) {
	fmt.Fprintf(l.cmd.OutOrStdout(), format, i...)
}
