package server

import "github.com/spf13/cobra"

// port http server port
var port int

// internal http server internal: v0 package
var internal bool

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Long:  "demo server",
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

func init() {
	Cmd.Flags().IntVar(&port, "port", 8080, "")
	Cmd.Flags().BoolVar(&internal, "internal", false, "")
}
