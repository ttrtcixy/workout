package workoutusecase

import (
	"context"
	"errors"

	"github.com/ttrtcixy/workout/internal/core/entities"
	"github.com/ttrtcixy/workout/internal/core/usecase/ports"
	apperrors "github.com/ttrtcixy/workout/internal/errors"
	"github.com/ttrtcixy/workout/internal/logger"
)

type CreateWorkoutUsecase struct {
	log  logger.Logger
	repo ports.CreateWorkoutRepository
}

func NewCreateWorkout(log logger.Logger, repo ports.CreateWorkoutRepository) *CreateWorkoutUsecase {
	return &CreateWorkoutUsecase{log: log, repo: repo}
}

func (u *CreateWorkoutUsecase) CreateWorkout(ctx context.Context, payload *entities.CreateWorkoutRequest) (response *entities.CreateWorkoutResponse, err error) {
	defer func() {
		if err != nil {
			var userErr apperrors.UserError
			if errors.As(err, &userErr) {
				return
			}
			
		}
	}()
	err = u.repo.CreateWorkout(ctx, payload)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
