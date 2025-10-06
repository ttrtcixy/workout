package handlers

import (
	"github.com/ttrtcixy/workout/internal/delivery/ports"
	"github.com/ttrtcixy/workout/internal/logger"
)

type Handlers struct {
	*CreateWorkout
}

func NewHandlers(log logger.Logger, usecase ports.UseCase) *Handlers {
	return &Handlers{
		NewCreateWorkout(log, usecase),
	}
}
