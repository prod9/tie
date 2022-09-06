package data

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chakrit/gendiff"
	"github.com/jmoiron/sqlx"
)

// language=PostgreSQL
const (
	CreateMigrationsTableSQL = `
		CREATE TABLE IF NOT EXISTS migrations
		(
			name     text PRIMARY KEY,
			up_sql   text NOT NULL,
			down_sql text NOT NULL
		);`
	ListMigrationsSQL = `
		SELECT *
		FROM migrations
		ORDER BY name ASC`
	UpdateMigrationSQL = `
		INSERT INTO migrations (name, up_sql, down_sql)
		VALUES ($1, $2, $3)
		ON CONFLICT (name) DO UPDATE
		SET up_sql = $2, down_sql = $3`
	PruneMigrationSQL = `
		DELETE FROM migrations
		WHERE name = $1`
)

const (
	MaxRollbacks = 1
	UpExt        = ".up.sql"
	DownExt      = ".down.sql"
)

type Action int

const (
	// ActionUpdate updates the migration content in the database
	ActionUpdate = Action(iota)
	// ActionIgnore ignores a new migration file that has yet to be run.
	ActionIgnore
	// ActionPrune prunes missing migration entries from the database
	ActionPrune
	// Exports the migration content from the database to a file
	// TODO: ActionExport
	//   This would be a different mode for `ActionPrune` so we'll need a switch
	//   for user to specify which action they want for migrations that is in
	//   the database but is missing from the filesystem. IMO This is a better
	//   default option than simply pruning as well. Also, if the switch is on,
	//   ActionUpdate should result in updating the filesystem from the DB rather
	//   than the other way around.

	// ActionMigrate runs all migrations that are yet to be ran
	ActionMigrate
	// ActionRollback rollbacks the most recent migration
	ActionRollback
)

type Intent int

const (
	IntentSync = Intent(iota)
	IntentMigrate
	IntentRollback
)

func (act Action) String() string {
	switch act {
	case ActionUpdate:
		return "update sql"
	case ActionIgnore:
		return "ignore"
	case ActionPrune:
		return "remove"
	case ActionMigrate:
		return "migrate"
	case ActionRollback:
		return "rollback"
	default:
		return "(unknown)"
	}

}

type Plan struct {
	Action    Action
	Migration Migration
}

func (p Plan) String() string {
	date, name := "", ""
	parts := strings.SplitN(p.Migration.Name, "_", 2)
	switch len(parts) {
	case 0:
		date, name = "209901010000", "(unknown)"
	case 1:
		date, name = parts[0], "(unknown)"
	default:
		date, name = parts[0], parts[1]
	}

	return fmt.Sprintf("%20s => %s %s", p.Action, date, name)
}

type Migration struct {
	Name    string `db:"name"`
	UpSQL   string `db:"up_sql"`
	DownSQL string `db:"down_sql"`
}

type migrationDiff struct {
	left  []Migration
	right []Migration
}

var _ gendiff.Interface = migrationDiff{}

func (d migrationDiff) LeftLen() int        { return len(d.left) }
func (d migrationDiff) RightLen() int       { return len(d.right) }
func (d migrationDiff) Equal(l, r int) bool { return d.left[l].Name == d.right[r].Name }

type Migrator struct {
	db  *sqlx.DB
	dir string
}

func NewMigrator(db *sqlx.DB, dir string) *Migrator {
	return &Migrator{db, dir}
}

func (m *Migrator) Plan(ctx context.Context, intent Intent) (actions []Plan, dirty bool, err error) {
	var (
		scope   Scope
		inFiles []Migration
		inDB    []Migration
	)

	if scope, err = NewScope(ctx, m.db); err != nil {
		return
	} else {
		defer scope.End(&err)
	}

	if inFiles, err = m.loadDir(scope.Context()); err != nil {
		return
	} else if inDB, err = m.loadFromDB(scope.Context()); err != nil {
		return
	}

	switch intent {
	case IntentSync:
		return m.planSync(inFiles, inDB)
	case IntentMigrate:
		return m.planMigrate(inFiles, inDB)
	case IntentRollback:
		return m.planRollback(inFiles, inDB)
	default:
		return nil, false, nil
	}
}

func (m *Migrator) planSync(inFiles []Migration, inDB []Migration) (actions []Plan, dirty bool, err error) {

	diffs := gendiff.Make(migrationDiff{inDB, inFiles})

	for _, d := range diffs {
		switch d.Op {
		case gendiff.Delete:
			dirty = true
			for lidx := d.Lstart; lidx < d.Lend; lidx++ {
				actions = append(actions, Plan{ActionPrune, inDB[lidx]})
			}

		case gendiff.Match:
			lidx, ridx := d.Lstart, d.Rstart
			for lidx < d.Lend && ridx < d.Rend {
				mDB, mFile := inDB[lidx], inFiles[ridx]
				if mDB.UpSQL != mFile.UpSQL || mDB.DownSQL != mFile.DownSQL {
					dirty = true
					actions = append(actions, Plan{ActionUpdate, mFile})
				}
				lidx += 1
				ridx += 1
			}
		}
	}

	return
}

func (m *Migrator) planMigrate(inFiles []Migration, inDB []Migration) (actions []Plan, dirty bool, err error) {
	diffs := gendiff.Make(migrationDiff{inDB, inFiles})

	for _, d := range diffs {
		switch d.Op {
		case gendiff.Insert:
			// some migrations were removed/changed prior to this migration, which means that
			// the db is likely not in the state that the migration expects it to be.
			if dirty {
				err = fmt.Errorf("db state divergence detected, please carefully review and re-sync")
				return
			}
			for ridx := d.Rstart; ridx < d.Rend; ridx++ {
				actions = append(actions, Plan{ActionMigrate, inFiles[ridx]})
			}

		case gendiff.Delete:
			dirty = true
			for lidx := d.Lstart; lidx < d.Lend; lidx++ {
				actions = append(actions, Plan{ActionPrune, inDB[lidx]})
			}

		case gendiff.Match:
			lidx, ridx := d.Lstart, d.Rstart
			for lidx < d.Lend && ridx < d.Rend {
				mDB, mFile := inDB[lidx], inFiles[ridx]
				if mDB.UpSQL != mFile.UpSQL || mDB.DownSQL != mFile.DownSQL {
					dirty = true
					actions = append(actions, Plan{ActionUpdate, mFile})
				}
				lidx += 1
				ridx += 1
			}
		}
	}

	return
}

func (m *Migrator) planRollback(inFiles []Migration, inDB []Migration) (actions []Plan, dirty bool, err error) {
	rollbackIdx := len(inDB) - MaxRollbacks

	diffs := gendiff.Make(migrationDiff{inDB, inFiles})
	for _, d := range diffs {
		switch d.Op {
		case gendiff.Insert:
			if d.Lstart <= rollbackIdx {
				err = fmt.Errorf("db state divergence detected, please carefully review and re-sync")
				return
			}

		case gendiff.Delete:
			dirty = true
			if d.Lstart <= rollbackIdx {
				err = fmt.Errorf("db state divergence detected, please carefully review and re-sync")
				return
			}
			for lidx := d.Lstart; lidx < d.Lend; lidx++ {
				actions = append(actions, Plan{ActionPrune, inDB[lidx]})
			}

		case gendiff.Match:
			lidx, ridx := d.Lstart, d.Rstart
			for lidx < d.Lend && ridx < d.Rend {
				mDB, mFile := inDB[lidx], inFiles[ridx]
				if mDB.UpSQL != mFile.UpSQL || mDB.DownSQL != mFile.DownSQL {
					dirty = true
					if d.Lstart <= rollbackIdx {
						err = fmt.Errorf("db state divergence detected, please carefully review and re-sync")
						return
					}
					actions = append(actions, Plan{ActionUpdate, mFile})
				}
				lidx += 1
				ridx += 1
			}
		}
	}

	// all prior states to the rollbacks seems OK, performing rollbacks (in reverse)
	for idx := len(inDB) - 1; idx >= rollbackIdx; idx-- {
		actions = append(actions, Plan{ActionRollback, inDB[idx]})
	}
	return
}

func (m *Migrator) Apply(ctx context.Context, plan Plan) (err error) {
	var (
		scope Scope
		mig   Migration
	)
	if scope, err = NewScope(ctx, m.db); err != nil {
		return
	} else {
		defer scope.End(&err)
	}

	mig = plan.Migration

	switch plan.Action {
	case ActionUpdate:
		if err = scope.Exec(UpdateMigrationSQL, mig.Name, mig.UpSQL, mig.DownSQL); err != nil {
			return
		}

	case ActionIgnore:
		// no-op

	case ActionPrune:
		if err = scope.Exec(PruneMigrationSQL, mig.Name); err != nil {
			return
		}

	case ActionMigrate:
		if err = scope.Exec(UpdateMigrationSQL, mig.Name, mig.UpSQL, mig.DownSQL); err != nil {
			return
		} else if err = scope.Exec(plan.Migration.UpSQL); err != nil {
			return
		}

	case ActionRollback:
		if err = scope.Exec(plan.Migration.DownSQL); err != nil {
			return
		} else if err = scope.Exec(PruneMigrationSQL, mig.Name); err != nil {
			return
		}

	default:
		err = fmt.Errorf("unknown Action in plan: %d", plan.Action)
	}

	return
}

func (m *Migrator) loadFromDB(ctx context.Context) (result []Migration, err error) {
	var scope Scope
	if scope, err = NewScope(ctx, m.db); err != nil {
		return nil, err
	} else {
		defer scope.End(&err)
	}

	if err = scope.Exec(CreateMigrationsTableSQL); err != nil {
		return
	} else if err = scope.Select(&result, ListMigrationsSQL); err != nil {
		return
	}

	return
}

func (m *Migrator) loadDir(context.Context) (result []Migration, err error) {
	files, err := filepath.Glob(filepath.Clean(m.dir) + "/*" + UpExt)
	if err != nil {
		return nil, err
	}

	var migration Migration
	sort.Sort(sort.StringSlice(files))
	for _, path := range files {
		downfile := path[:len(path)-len(UpExt)] + DownExt
		basename := filepath.Base(path)
		migration.Name = basename[:len(basename)-len(UpExt)]

		if _, err := os.Stat(downfile); os.IsNotExist(err) {
			return nil, fmt.Errorf("missing down migration: %s", downfile)
		}

		if bytes, err := ioutil.ReadFile(path); err != nil {
			return nil, fmt.Errorf("i/o: %w", err)
		} else {
			migration.UpSQL = string(bytes)
		}

		if bytes, err := ioutil.ReadFile(downfile); err != nil {
			return nil, fmt.Errorf("i/o: %w", err)
		} else {
			migration.DownSQL = string(bytes)
		}

		result = append(result, migration)
	}

	return
}
