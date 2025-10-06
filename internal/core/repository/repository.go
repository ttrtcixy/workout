package repository

import (
	"context"
	"github.com/ttrtcixy/workout/internal/logger"
	storage "github.com/ttrtcixy/workout/internal/storage/pg"
)

type Repository struct {
	log logger.Logger
	DB  storage.DB
}

func NewRepository(ctx context.Context, log logger.Logger, db storage.DB) *Repository {
	return &Repository{
		log,
		db,
	}
}

func (r *Repository) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	err := r.DB.RunInTx(ctx, fn)
	if err != nil {
		return err
	}
	return nil
}
