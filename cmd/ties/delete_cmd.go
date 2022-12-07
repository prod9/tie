package ties

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"tie.prodigy9.co/client"
	"tie.prodigy9.co/cmd/prompts"
	"tie.prodigy9.co/config"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete slug",
	Short: "Delete a tie at /slug",
	RunE:  runDeleteCmd,
}

func runDeleteCmd(cmd *cobra.Command, args []string) error {
	cfg := config.MustConfigure()
	p := prompts.New(cfg, args)

	slug := p.Str("slug")

	client, err := client.NewClient(config.MustConfigure())
	if err != nil {
		return err
	}

	tie, err := client.DeleteTie(slug)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, tie)
	return nil
}
