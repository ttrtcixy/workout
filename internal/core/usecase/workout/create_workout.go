package workoutusecase

import (
	"context"
	"github.com/ttrtcixy/workout/internal/core/entities"
	"github.com/ttrtcixy/workout/internal/core/usecase/ports"
	"github.com/ttrtcixy/workout/internal/logger"
)

type CreateWorkoutUsecase struct {
	log  logger.Logger
	repo ports.CreateWorkoutRepository
}

func NewCreateWorkout(log logger.Logger, repo ports.CreateWorkoutRepository) *CreateWorkoutUsecase {
	return &CreateWorkoutUsecase{log: log, repo: repo}
}

func (u *CreateWorkoutUsecase) CreateWorkout(ctx context.Context, payload *entities.CreateWorkoutRequest) (*entities.CreateWorkoutResponse, error) {
	err := u.repo.CreateWorkout(ctx, payload)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
