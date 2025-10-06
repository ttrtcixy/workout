package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ttrtcixy/workout/internal/config"
	"github.com/ttrtcixy/workout/internal/logger"
)

type contextKey string

var txCtxKey = contextKey("tx")

type db struct {
	cfg  *config.DB
	pool *pgxpool.Pool
	log  logger.Logger
}

// todo check
func (db *db) RunInTx(ctx context.Context, fn func(context.Context) error) (err error) {
	// будут ли проблемы из-за того что не указатель на транзакцию
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				db.log.Warning(rollbackErr.Error())
			}
		}
	}()

	// будут ли проблемы из-за
	ctx = withValue(ctx, tx)

	if err := fn(ctx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		db.log.Warning(err.Error())
		return err
	}
	return nil
}

// todo change to Tx
func withValue(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey, tx)
}

func value(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txCtxKey).(pgx.Tx)
	return tx, ok
}

func (db *db) Close(ctx context.Context) error {
	const op = "db.close"
	return func() error {
		db.pool.Close()
		return nil
	}()
}

func (db *db) Exec(ctx context.Context, query Query) (CommandTag, error) {
	db.logQuery(query)

	if val, ok := value(ctx); ok {

		return val.Exec(ctx, query.Query(), query.Args()...)
	}

	return db.pool.Exec(ctx, query.Query(), query.Args()...)
}
func (db *db) Query(ctx context.Context, query Query) (Rows, error) {
	db.logQuery(query)
	if val, ok := value(ctx); ok {
		rw, err := val.Query(ctx, query.Query(), query.Args()...)
		if err != nil {
			return nil, err
		}
		return &rows{rw}, nil
	}

	rw, err := db.pool.Query(ctx, query.Query(), query.Args()...)
	if err != nil {
		return nil, err
	}

	//return db.connect.Query(ctx, query.Query(), query.Args()...)
	return &rows{rw}, nil
}
func (db *db) QueryRow(ctx context.Context, query Query) Row {
	db.logQuery(query)
	if val, ok := value(ctx); ok {
		return val.QueryRow(ctx, query.Query(), query.Args()...)
	}

	return db.pool.QueryRow(ctx, query.Query(), query.Args()...)
}

func New(ctx context.Context, log logger.Logger, cfg *config.DB) (DB, error) {
	const op = "storage.New"
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("op: %s, error connect to database: %s", op, err.Error())
	}
	storage := &db{
		cfg:  cfg,
		pool: pool,
		log:  log,
	}

	if err := storage.pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("op: %s, error connect to database: %s", op, err.Error())
	}

	return storage, nil
}

func (db *db) logQuery(query Query) {
	db.log.Debug(query.String())
}
