package cli

import (
	"github.com/cpekyaman/devtools/cli/cmd/codegen"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "gdevl",
		Short: "goits developer tools cli",
		Long:  "command line interface to perform certain devl tasks",
	}

	rootCmd.AddCommand(codegen.CodegenCommand)
}

// Execute runs the root cli command as the main entry point.
func Execute() error {
	return rootCmd.Execute()
}
