package data

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"tie.prodigy9.co/cmd/prompts"
	"tie.prodigy9.co/config"
	"tie.prodigy9.co/internal"

	"github.com/gobuffalo/flect"
	"github.com/spf13/cobra"
)

const upMigrationTemplate = `-- vim: filetype=SQL
CREATE TABLE dummy (
	id TEXT PRIMARY KEY,
	ctime TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`
const downMigrationTemplate = `-- vim: filetype=SQL
DROP TABLE dummy;
`

var newMigrationCmd = &cobra.Command{
	Use:   "new-migration (name)",
	Short: "Creates a new migration file with timestamps and the given name",
	RunE:  runNewMigrationCmd,
}

func runNewMigrationCmd(cmd *cobra.Command, args []string) (err error) {
	defer internal.WrapErr("new-migration", &err)

	cfg := config.MustConfigure()
	prompt := prompts.New(cfg, args)
	name := prompt.Str("name of migration")

	name = time.Now().Format("200601021504") + "_" + flect.Underscore(name)
	upname, downname := name+".up.sql", name+".down.sql"

	uppath, err := filepath.Abs(filepath.Join("./data/middlewares", upname))
	if err != nil {
		return err
	}
	downpath, err := filepath.Abs(filepath.Join("./data/middlewares", downname))
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, uppath)
	fmt.Fprintln(os.Stdout, downpath)
	if !prompt.YesNo("create these files") {
		log.Fatalln("aborted")
	}

	if err := ioutil.WriteFile(uppath, []byte(upMigrationTemplate), 0644); err != nil {
		return err
	} else if err := ioutil.WriteFile(downpath, []byte(downMigrationTemplate), 0644); err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if strings.TrimSpace(editor) == "" {
		editor = "/usr/bin/vi"
	}

	proc := exec.Command(editor, uppath, downpath)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	return proc.Run()
}
