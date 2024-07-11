package repo

import (
	"workout/models"

	"gorm.io/gorm"
)

type WorkoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) *WorkoutRepository {
	return &WorkoutRepository{db}
}

func (r *WorkoutRepository) CreateWorkout(workout *models.Workout) error {
	return r.db.Create(workout).Error
}

func (r *WorkoutRepository) GetWorkoutByID(id int) (*models.Workout, error) {
	var workout models.Workout
	err := r.db.First(&workout, id).Error
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (r *WorkoutRepository) UpdateWorkout(workout *models.Workout, updates map[string]interface{}) error {
	return r.db.Model(workout).Updates(updates).Error
}

func (r *WorkoutRepository) DeleteWorkout(id int) error {
	return r.db.Delete(&models.Workout{}, id).Error
}

func (r *WorkoutRepository) GetWorkoutDetails(id int) (*models.WorkoutDetails, error) {
	var workoutDetails models.WorkoutDetails

	var workout models.Workout
	err := r.db.Preload("Exercises.Sets").Where("id = ?", id).First(&workout).Error
	if err != nil {
		return nil, err
	}

	workoutDetails.ID = workout.ID
	workoutDetails.Name = workout.Name

	for _, we := range workout.Exercises {
		var exercise models.Exercise
		r.db.First(&exercise, we.ExerciseID)

		exerciseDetails := models.ExerciseDetails{
			ID:   exercise.ID,
			Name: exercise.Name,
			Sets: we.Sets,
		}

		workoutDetails.Exercises = append(workoutDetails.Exercises, exerciseDetails)
	}

	return &workoutDetails, nil
}
