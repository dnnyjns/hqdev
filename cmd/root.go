package cmd

import (
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	var (
		rootCmd = &cobra.Command{
			Use:   "hqdev",
			Short: "hqdev, a CLI to help with hq development",
		}
	)

	cobra.OnInitialize(RootConfig.initConfig)

	// Register Command
	rootCmd.AddCommand(generateRestoreCommand())
	rootCmd.AddCommand(generateResetCommand())

	// Run Command
	err := rootCmd.Execute()
	onError(err)
}
