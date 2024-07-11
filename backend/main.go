package main

import (
	"fmt"
	"log"
	"os"

	"workout/handlers"
	"workout/models"
	"workout/repo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
		port     = os.Getenv("DB_PORT")
		sslmode  = os.Getenv("DB_SSLMODE")
	)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Exercise{})
	db.AutoMigrate(&models.Workout{})
	db.AutoMigrate(&models.WorkoutExercise{})
	db.AutoMigrate(&models.Set{})
}

func main() {
	fmt.Println("Connected to PostgreSQL!")

	// repos
	userRepo := repo.NewUserRepository(db)
	exerciseRepo := repo.NewExerciseRepository(db)
	workoutRepo := repo.NewWorkoutRepository(db)
	setRepo := repo.NewSetRepository(db)
	workoutExerciseRepo := repo.NewWorkoutExerciseRepository(db)

	// handlers
	userHandler := handlers.NewUserHandler(userRepo)
	exerciseHandler := handlers.NewExerciseHandler(exerciseRepo)
	workoutHandler := handlers.NewWorkoutHandler(workoutRepo, exerciseRepo)
	setHandler := handlers.NewSetHandler(setRepo, workoutExerciseRepo)

	r := gin.Default()

	// Enable CORS for the router
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// user
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	// exercise
	r.POST("/exercises", exerciseHandler.CreateExercise)
	r.GET("/exercises/:id", exerciseHandler.GetExerciseByID)
	r.PUT("/exercises/:id", exerciseHandler.UpdateExercise)
	r.DELETE("/exercises/:id", exerciseHandler.DeleteExercise)

	// workout
	r.POST("/workouts", workoutHandler.CreateWorkout)
	r.POST("/workouts/:id/exercises", workoutHandler.AddExerciseToWorkout)
	r.GET("/workouts/:id", workoutHandler.GetWorkoutDetails)
	r.PUT("/workouts/:id", workoutHandler.UpdateWorkout)
	r.DELETE("/workouts/:id", workoutHandler.DeleteWorkout)

	// sets - exercise_id = workoutExerciseId
	r.POST("/workouts/:id/exercises/:exercise_id/sets", setHandler.AddSetToExercise)
	r.GET("/workouts/:id/exercises/:exercise_id/sets", setHandler.GetSetsForExercise)

	r.Run()
}
