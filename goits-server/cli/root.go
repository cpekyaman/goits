package cli

import (
	"github.com/cpekyaman/goits/cli/cmd"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "goits",
		Short: "goits server cli",
		Long:  "command line interface to manage goits server",
	}

	rootCmd.AddCommand(cmd.MigrateCommand)
	rootCmd.AddCommand(cmd.ServerCommand)
}

// Execute runs the root cli command as the main entry point.
func Execute() error {
	return rootCmd.Execute()
}
