package ties

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "ties",
	Short: "Work with ties as a client",
}

func init() {
	Cmd.AddCommand(
		CreateCmd,
		ListCmd,
	)
}
