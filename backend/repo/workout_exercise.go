package repo

import (
	"workout/models"

	"gorm.io/gorm"
)

type WorkoutExerciseRepository struct {
	db *gorm.DB
}

func NewWorkoutExerciseRepository(db *gorm.DB) *WorkoutExerciseRepository {
	return &WorkoutExerciseRepository{db}
}

func (r *WorkoutExerciseRepository) GetWorkoutExerciseByID(id int) (*models.WorkoutExercise, error) {
	var exercise models.WorkoutExercise
	err := r.db.First(&exercise, id).Error
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}
