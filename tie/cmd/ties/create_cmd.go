package ties

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"tie.prodigy9.co/client"
	"tie.prodigy9.co/cmd/prompts"
	"tie.prodigy9.co/config"
)

var CreateCmd = &cobra.Command{
	Use:   "create slug url",
	Short: "Create a new tie at /slug pointing to url",
	RunE:  runCreateCmd,
}

func runCreateCmd(cmd *cobra.Command, args []string) error {
	cfg := config.MustConfigure()
	p := prompts.New(cfg, args)

	slug := p.Str("slug")
	target := p.Str("target url")

	client, err := client.NewClient(config.MustConfigure())
	if err != nil {
		return err
	}

	tie, err := client.CreateTie(slug, target)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, tie)
	return nil
}
