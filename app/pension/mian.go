package pension

import "github.com/spf13/cobra"

func main() {
	cmd := &cobra.Command{
		Use:   "pension",
		Short: "pension",
		Long:  "a project named pension",
	}

	cmd.Execute()
}
