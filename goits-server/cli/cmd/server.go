package cmd

import (
	"fmt"

	"github.com/cpekyaman/goits/server"
	"github.com/spf13/cobra"
)

var ServerCommand *cobra.Command

func init() {
	ServerCommand = &cobra.Command{
		Use:     "server",
		Aliases: []string{"svc"},
		Short:   "manage goits server",
		Long:    "server management commands such as start and stop",
	}

	start := &cobra.Command{
		Use:   "start",
		Short: "start the server",
		Long:  "start the server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start()
		},
	}

	stop := &cobra.Command{
		Use:   "stop",
		Short: "stop the server",
		Long:  "stop the server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("stop server")
		},
	}

	ServerCommand.AddCommand(start, stop)
}
