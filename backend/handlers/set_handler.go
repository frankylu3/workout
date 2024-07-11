package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"workout/models"
	"workout/repo"

	"github.com/gin-gonic/gin"
)

type SetHandler struct {
	SetRepository      *repo.SetRepository
	ExerciseRepository *repo.WorkoutExerciseRepository
}

func NewSetHandler(setRepo *repo.SetRepository, exerciseRepo *repo.WorkoutExerciseRepository) *SetHandler {
	return &SetHandler{setRepo, exerciseRepo}
}

func (h *SetHandler) AddSetToExercise(c *gin.Context) {
	var set models.Set
	err := c.ShouldBindJSON(&set)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if set.SetNumber < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "set number must be greater than 0"})
	}

	exerciseId, err := strconv.Atoi(c.Param("exercise_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exercise ID"})
		return
	}
	sets, err := h.SetRepository.GetSetsForExercise(exerciseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, s := range sets {
		if s.SetNumber == set.SetNumber {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%d already exists", set.SetNumber)})
			return
		}
	}

	_, err = h.ExerciseRepository.GetWorkoutExerciseByID(exerciseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	set.WorkoutExerciseID = uint(exerciseId)
	err = h.SetRepository.CreateSet(&set)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, set)
}

func (h *SetHandler) GetSetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid set ID"})
		return
	}

	set, err := h.SetRepository.GetSetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, set)
}

func (h *SetHandler) DeleteSet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid set ID"})
		return
	}

	err = h.SetRepository.DeleteSet(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Set deleted"})
}

func (h *SetHandler) GetSetsForExercise(c *gin.Context) {
	exerciseId, err := strconv.Atoi(c.Param("exercise_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid set ID"})
		return
	}

	sets, err := h.SetRepository.GetSetsForExercise(exerciseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sets)
}
