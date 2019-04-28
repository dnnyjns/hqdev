package cmd

import (
	"github.com/spf13/cobra"
)

func generateResetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset saved configuration",
		Run: func(cmd *cobra.Command, args []string) {
			RootConfig.Reset()
		},
	}
}
