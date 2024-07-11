package models

import "gorm.io/gorm"

// Represents the many-to-many relationships between workouts and exercises
type WorkoutExercise struct {
	gorm.Model
	WorkoutID  uint `gorm:"not null"`
	ExerciseID uint `gorm:"not null"`
	Sets       []Set
}
