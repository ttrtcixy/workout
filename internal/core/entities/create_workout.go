package entities

import "time"

type CreateWorkoutRequest struct {
	WorkoutName   string
	WorkoutType   string
	ExerciseCount int
	Exercises     []*Exercise

	UserID int
}

type Exercise struct {
	Name                string
	MuscleGroup         string
	RestTime            time.Duration
	ExecutionType       string
	ExercisesCombining  string
	NumberOfApproaches  int16
	NumberOfRepetitions int16
}

type CreateWorkoutResponse struct {
}
