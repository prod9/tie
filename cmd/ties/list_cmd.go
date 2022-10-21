package ties

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"tie.prodigy9.co/client"
	"tie.prodigy9.co/config"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ties",
	RunE:  runListCmd,
}

func runListCmd(cmd *cobra.Command, args []string) error {
	client, err := client.NewClient(config.MustConfigure())
	if err != nil {
		return err
	}

	ties, err := client.GetTies()
	if err != nil {
		return err
	}

	for _, tie := range ties {
		fmt.Fprintln(os.Stdout, tie)
	}
	return nil
}
