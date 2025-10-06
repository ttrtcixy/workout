package usecase

import (
	"github.com/ttrtcixy/workout/internal/core/usecase/ports"
	workoutusecase "github.com/ttrtcixy/workout/internal/core/usecase/workout"
	"github.com/ttrtcixy/workout/internal/logger"
)

type UseCase struct {
	*WorkoutUsecase
}

func NewUseCase(log logger.Logger, repo ports.Repository) *UseCase {
	return &UseCase{
		NewWorkoutUsecase(log, repo),
	}
}

type WorkoutUsecase struct {
	*workoutusecase.CreateWorkoutUsecase
}

func NewWorkoutUsecase(log logger.Logger, repo ports.Repository) *WorkoutUsecase {
	return &WorkoutUsecase{
		workoutusecase.NewCreateWorkout(log, repo),
	}
}
