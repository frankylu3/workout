package repo

import (
	"workout/models"

	"gorm.io/gorm"
)

type ExerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) *ExerciseRepository {
	return &ExerciseRepository{db}
}

func (r *ExerciseRepository) CreateExercise(exercise *models.Exercise) error {
	return r.db.Create(exercise).Error
}

func (r *ExerciseRepository) GetExerciseByID(id int) (*models.Exercise, error) {
	var exercise models.Exercise
	err := r.db.First(&exercise, id).Error
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (r *ExerciseRepository) UpdateExercise(exercise *models.Exercise, updates map[string]interface{}) error {
	return r.db.Model(exercise).Updates(updates).Error
}

func (r *ExerciseRepository) DeleteExercise(id int) error {
	return r.db.Delete(&models.Exercise{}, id).Error
}

func (r *ExerciseRepository) AddExerciseToWorkout(exercise *models.WorkoutExercise) error {
	return r.db.Create(exercise).Error
}
