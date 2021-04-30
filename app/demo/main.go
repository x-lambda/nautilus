package main

import (
	"github.com/x-lambda/nautilus/app/demo/cmd/help"
	"github.com/x-lambda/nautilus/app/demo/cmd/job"
	"github.com/x-lambda/nautilus/app/demo/cmd/server"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "",
		Short: "protoc-gen-gin-example",
		Long:  "protoc-gen-gin example",
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
