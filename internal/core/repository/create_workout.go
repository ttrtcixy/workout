package repository

import (
	"context"
	"fmt"
	"github.com/ttrtcixy/workout/internal/core/entities"
	"github.com/ttrtcixy/workout/internal/core/repository/query"
	"strings"
)

const (
	createWorkoutSQL = `
		insert into workouts (workout_name, workout_type, exercise_count, user_id)
		values ($1, $2, $3, $4)
		returning workout_id
`
	createExercise = `
		INSERT INTO exercises (
			exercise_name, 
		    muscle_group, 
			execution_type, 
			exercises_combining,
			number_of_repetitions, 
			number_of_approaches, 
		    rest_time,
			workout_id, 
			user_id) VALUES %s
`
)

func (r *Repository) CreateWorkout(ctx context.Context, payload *entities.CreateWorkoutRequest) error {
	err := r.RunInTx(ctx, func(ctx context.Context) error {
		workoutID, err := r.createWorkout(ctx, payload)
		if err != nil {
			return err
		}

		err = r.createExercise(ctx, payload, workoutID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) createWorkout(ctx context.Context, payload *entities.CreateWorkoutRequest) (workoutID int, err error) {
	q := &query.Query{
		Name:      "create workout for user",
		RawQuery:  createWorkoutSQL,
		Arguments: []any{payload.WorkoutName, payload.WorkoutType, payload.ExerciseCount, payload.UserID},
	}

	err = r.DB.QueryRow(ctx, q).Scan(&workoutID)
	if err != nil {
		return 0, err
	}

	return workoutID, nil
}

func (r *Repository) createExercise(ctx context.Context, payload *entities.CreateWorkoutRequest, workoutID int) error {
	var valuesBuilder strings.Builder

	for i, ex := range payload.Exercises {
		if i > 0 {
			valuesBuilder.WriteString(",")
		}
		valuesBuilder.WriteString(fmt.Sprintf(
			"('%s', '%s', '%s', '%s', %d, %d, '%s', %d, %d)",
			ex.Name,
			ex.MuscleGroup,
			ex.ExecutionType,
			ex.ExercisesCombining,
			ex.NumberOfRepetitions,
			ex.NumberOfApproaches,
			ex.RestTime.String(),
			workoutID,
			payload.UserID,
		))
	}

	q := &query.Query{
		Name:      "add exercises",
		RawQuery:  fmt.Sprintf(createExercise, valuesBuilder.String()),
		Arguments: nil,
	}

	_, err := r.DB.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}
