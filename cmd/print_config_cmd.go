package cmd

import (
	"fmt"
	"os"

	"tie.prodigy9.co/config"

	"github.com/spf13/cobra"
)

var printConfigCmd = &cobra.Command{
	Use:   "print-config",
	Short: "Prints current effective configuration.",
	RunE:  runPrintConfigCmd,
}

func runPrintConfigCmd(cmd *cobra.Command, args []string) error {
	configs := config.MustConfigure().AllConfigurations()
	for key, value := range configs {
		fmt.Fprintln(os.Stdout, key, "=", value)
	}
	return nil
}
