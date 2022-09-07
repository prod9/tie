package data

import (
	"tie.prodigy9.co/data"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback [middlewares-dir]",
	Short: "Revert one previously ran migration.",
	RunE:  runRollbackCmd,
}

func runRollbackCmd(cmd *cobra.Command, args []string) error {
	return runMigration(data.IntentRollback, args)
}
