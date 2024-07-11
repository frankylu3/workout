package repo

import (
	"workout/models"

	"gorm.io/gorm"
)

type SetRepository struct {
	db *gorm.DB
}

func NewSetRepository(db *gorm.DB) *SetRepository {
	return &SetRepository{db}
}

func (r *SetRepository) CreateSet(set *models.Set) error {
	return r.db.Create(set).Error
}

func (r *SetRepository) GetSetByID(id int) (*models.Set, error) {
	var set models.Set
	err := r.db.First(&set, id).Error
	if err != nil {
		return nil, err
	}
	return &set, nil
}

func (r *SetRepository) UpdateSet(set *models.Set, updates map[string]interface{}) error {
	return r.db.Model(set).Updates(updates).Error
}

func (r *SetRepository) DeleteSet(id int) error {
	return r.db.Delete(&models.Set{}, id).Error
}

func (r *SetRepository) GetSetsForExercise(exerciseId int) (sets []models.Set, err error) {
	err = r.db.Where("workout_exercise_id", exerciseId).Find(&sets).Error
	if err != nil {
		return nil, err
	}

	return sets, nil
}
