package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ttrtcixy/workout/internal/core/entities"
	"github.com/ttrtcixy/workout/internal/delivery/http/dtos"
	"github.com/ttrtcixy/workout/internal/delivery/ports"
	"github.com/ttrtcixy/workout/internal/logger"
)

type CreateWorkout struct {
	log     logger.Logger
	usecase ports.CreateWorkoutUsecase
}

func NewCreateWorkout(log logger.Logger, usecase ports.CreateWorkoutUsecase) *CreateWorkout {
	return &CreateWorkout{log: log, usecase: usecase}
}

func (h *CreateWorkout) Run(c *gin.Context) {
	payload := &dtos.CreateWorkoutRequest{}

	err := c.ShouldBindJSON(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный json",
		})
		return
	}

	err = payload.Validate()
	if err != nil {
		var vErr *dtos.ValidationErrors
		if errors.As(err, &vErr) {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	// todo user id from ctx
	userID := 0

	result, err := h.usecase.CreateWorkout(c.Request.Context(), h.dtoToEntity(payload, userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_ = result
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *CreateWorkout) dtoToEntity(payload *dtos.CreateWorkoutRequest, userID int) *entities.CreateWorkoutRequest {
	es := make([]*entities.Exercise, 0, len(payload.Exercises))
	for _, v := range payload.Exercises {
		t, _ := time.ParseDuration(v.RestTime)
		es = append(es, &entities.Exercise{
			Name:                v.Name,
			MuscleGroup:         v.MuscleGroup,
			RestTime:            t,
			ExecutionType:       v.ExecutionType,
			ExercisesCombining:  v.ExercisesCombining,
			NumberOfApproaches:  v.NumberOfApproaches,
			NumberOfRepetitions: v.NumberOfRepetitions,
		})
	}
	return &entities.CreateWorkoutRequest{
		WorkoutName:   payload.WorkoutName,
		WorkoutType:   payload.WorkoutType,
		ExerciseCount: len(payload.Exercises),
		Exercises:     es,
		UserID:        userID,
	}
}
