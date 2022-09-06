package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"tie.prodigy9.co/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type dbContextKey struct{}

func FromContext(ctx context.Context) *sqlx.DB {
	return ctx.Value(dbContextKey{}).(*sqlx.DB)
}

func NewContext(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, dbContextKey{}, db)
}

func MustConnect(cfg *config.Config) *sqlx.DB {
	if db, err := Connect(cfg); err != nil {
		log.Panicln(err)
		return nil
	} else {
		return db
	}
}

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	if db, err := sqlx.Open("postgres", cfg.DatabaseURL()); err != nil {
		return nil, fmt.Errorf("database: %w", err)
	} else {
		return db, nil
	}
}

func NewScope(ctx context.Context, db *sqlx.DB) (Scope, error) {
	if db == nil {
		db = FromContext(ctx)
	}

	if impl, err := newScope(ctx, db); err != nil {
		return nil, err
	} else {
		return impl, nil
	}
}

func IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func Exec(ctx context.Context, sql string, args ...any) error {
	_, err := Run(ctx, func(s Scope) (any, error) {
		return nil, s.Exec(sql, args...)
	})
	return err
}

func Run[T any](ctx context.Context, action func(s Scope) (T, error)) (result T, err error) {
	var scope Scope
	if scope, err = NewScope(ctx, nil); err != nil {
		return
	} else {
		defer scope.End(&err)
		return action(scope)
	}
}
