package models

type WorkoutDetails struct {
	ID        uint   // workout id
	Name      string // workout name
	Exercises []ExerciseDetails
}

type ExerciseDetails struct {
	ID   uint   // exercise workout id
	Name string // exercise name
	Sets []Set
}
