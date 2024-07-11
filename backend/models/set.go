package models

import "gorm.io/gorm"

// represents a set for a workout exercise
type Set struct {
	gorm.Model
	WorkoutExerciseID uint
	SetNumber         int `gorm:"not null"`
	Reps              int
	Weight            float64
}
