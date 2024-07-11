package models

import "gorm.io/gorm"

// a workout belongs to one user
type Workout struct {
	gorm.Model
	Name      string `gorm:"not null"`
	UserID    uint
	Exercises []WorkoutExercise
}
