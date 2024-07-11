package handlers

import (
	"net/http"
	"strconv"
	"workout/models"
	"workout/repo"

	"github.com/gin-gonic/gin"
)

type WorkoutHandler struct {
	WorkoutRepository  *repo.WorkoutRepository
	ExerciseRepository *repo.ExerciseRepository
}

func NewWorkoutHandler(workoutRepo *repo.WorkoutRepository, exerciseRepo *repo.ExerciseRepository) *WorkoutHandler {
	return &WorkoutHandler{workoutRepo, exerciseRepo}
}

func (h *WorkoutHandler) CreateWorkout(c *gin.Context) {
	var workout models.Workout
	err := c.ShouldBindJSON(&workout)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if workout.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be empty"})
		return
	}

	err = h.WorkoutRepository.CreateWorkout(&workout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workout)
}

func (h *WorkoutHandler) GetWorkoutByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout ID"})
		return
	}

	workout, err := h.WorkoutRepository.GetWorkoutByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workout)
}

func (h *WorkoutHandler) UpdateWorkout(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout ID"})
		return
	}

	workout, err := h.WorkoutRepository.GetWorkoutByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var updates map[string]interface{}
	err = c.ShouldBindJSON(&updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.WorkoutRepository.UpdateWorkout(workout, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workout)
}

func (h *WorkoutHandler) DeleteWorkout(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout ID"})
		return
	}

	if err := h.WorkoutRepository.DeleteWorkout(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workout deleted"})
}

func (h *WorkoutHandler) AddExerciseToWorkout(c *gin.Context) {
	// id = which workout to add to
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout ID"})
		return
	}

	var workoutExercise models.WorkoutExercise
	err = c.ShouldBindJSON(&workoutExercise)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check that exercise and workout exists
	_, err = h.WorkoutRepository.GetWorkoutByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if workoutExercise.ExerciseID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exercise ID"})
		return
	}
	_, err = h.ExerciseRepository.GetExerciseByID(int(workoutExercise.ExerciseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	workoutExercise.WorkoutID = uint(id)
	err = h.ExerciseRepository.AddExerciseToWorkout(&workoutExercise)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise added to workout"})
}

func (h *WorkoutHandler) GetWorkoutDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout ID"})
		return
	}

	workoutDetails, err := h.WorkoutRepository.GetWorkoutDetails(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workoutDetails)
}
