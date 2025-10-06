package ports

import (
	"context"
	"github.com/ttrtcixy/workout/internal/core/entities"
)

type Repository interface {
	WorkoutRepository
}

type Tx interface {
	RunInTx(ctx context.Context, fn func(context.Context) error) error
}

type WorkoutRepository interface {
	CreateWorkoutRepository
}

type CreateWorkoutRepository interface {
	CreateWorkout(ctx context.Context, payload *entities.CreateWorkoutRequest) error
	Tx
}
