package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type DB interface {
	SQL
	Transaction
}

type SQL interface {
	Exec(ctx context.Context, query Query) (CommandTag, error)
	Query(ctx context.Context, query Query) (Rows, error)
	QueryRow(ctx context.Context, query Query) Row
	Close(ctx context.Context) error
}

type Transaction interface {
	RunInTx(ctx context.Context, fn func(context.Context) error) error
}

type CommandTag interface {
	RowsAffected() int64
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
	Err() error
	CommandTag() CommandTag
	Values() ([]any, error)
	RawValues() [][]byte
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Query interface {
	QueryName() string
	Query() string
	Args() []any
	String() string
}

type rows struct {
	pgx.Rows
}

func (r *rows) CommandTag() CommandTag {
	return r.Rows.CommandTag()
}
