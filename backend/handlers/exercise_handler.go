package handlers

import (
	"net/http"
	"strconv"
	"workout/models"
	"workout/repo"

	"github.com/gin-gonic/gin"
)

type ExerciseHandler struct {
	ExerciseRepository *repo.ExerciseRepository
}

func NewExerciseHandler(exerciseRepo *repo.ExerciseRepository) *ExerciseHandler {
	return &ExerciseHandler{exerciseRepo}
}

func (h *ExerciseHandler) CreateExercise(c *gin.Context) {
	var exercise models.Exercise
	err := c.ShouldBindJSON(&exercise)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if exercise.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be empty"})
		return
	}

	err = h.ExerciseRepository.CreateExercise(&exercise)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

func (h *ExerciseHandler) GetExerciseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exercise ID"})
		return
	}

	exercise, err := h.ExerciseRepository.GetExerciseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func (h *ExerciseHandler) UpdateExercise(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exercise ID"})
		return
	}

	exercise, err := h.ExerciseRepository.GetExerciseByID(id)
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

	err = h.ExerciseRepository.UpdateExercise(exercise, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func (h *ExerciseHandler) DeleteExercise(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exercise ID"})
		return
	}

	if err := h.ExerciseRepository.DeleteExercise(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted"})
}
