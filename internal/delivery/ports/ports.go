package ports

import (
	"context"
	"github.com/ttrtcixy/workout/internal/core/entities"
)

type UseCase interface {
	CreateWorkoutUsecase
}

type CreateWorkoutUsecase interface {
	CreateWorkout(ctx context.Context, payload *entities.CreateWorkoutRequest) (*entities.CreateWorkoutResponse, error)
}
