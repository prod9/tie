package cmd

import (
	"github.com/spf13/cobra"
	"tie.prodigy9.co/cmd/data"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts the API application",
}

func init() {
	rootCmd.AddCommand(
		data.Cmd,
		printConfigCmd,
		serveCmd,
	)
}
