package main

import (
	"nautilus/app/demo/cmd/help"
	"nautilus/app/demo/cmd/job"
	"nautilus/app/demo/cmd/server"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "",
		Short: "demo",
		Long:  "demo",
		Run: func(cmd *cobra.Command, args []string) {
			help.Usage()
		},
	}

	rootCmd.AddCommand(
		help.Cmd,
		server.Cmd,
		job.Cmd,
	)

	rootCmd.Execute()
}
