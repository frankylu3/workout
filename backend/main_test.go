package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"workout/handlers"
	"workout/models"
	"workout/repo"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testDB              *gorm.DB
	testUserRepo        *repo.UserRepository
	testUserHandler     *handlers.UserHandler
	testExerciseRepo    *repo.ExerciseRepository
	testExerciseHandler *handlers.ExerciseHandler
)

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME_TEST")
		port     = os.Getenv("DB_PORT")
		sslmode  = os.Getenv("DB_SSLMODE")
	)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	testDB = db
	testDB.AutoMigrate(&models.User{})
	testDB.AutoMigrate(&models.Exercise{})
	testDB.AutoMigrate(&models.Workout{})
	testDB.AutoMigrate(&models.WorkoutExercise{})
	testDB.AutoMigrate(&models.Set{})

	code := m.Run()
	os.Exit(code)
}

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	tx := testDB.Begin()

	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	return tx, func() {
		tx.Rollback()
	}
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Initialize repositories and handlers
	testUserRepo = repo.NewUserRepository(db)
	testUserHandler = handlers.NewUserHandler(testUserRepo)
	testExerciseRepo = repo.NewExerciseRepository(db)
	testExerciseHandler = handlers.NewExerciseHandler(testExerciseRepo)

	// Define routes for testing
	r.POST("/exercises", testExerciseHandler.CreateExercise)
	r.GET("/exercises/:id", testExerciseHandler.GetExerciseByID)
	r.PUT("/exercises/:id", testExerciseHandler.UpdateExercise)
	r.DELETE("/exercises/:id", testExerciseHandler.DeleteExercise)

	return r
}

func TestCreateExercise(t *testing.T) {
	Convey("Given a database and an exercise", t, func() {
		db, cleanup := setupTestDB(t)
		defer cleanup()
		exercise := models.Exercise{
			Name: "Bench Press",
		}
		w := httptest.NewRecorder()
		r := setupRouter(db)
		Convey("When a create request is sent", func() {
			Convey("And everything is given", func() {
				body, _ := json.Marshal(exercise)
				req, _ := http.NewRequest("POST", "/exercises", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				Convey("Then the exercise is created", func() {
					r.ServeHTTP(w, req)
					var createdExercise models.Exercise
					json.Unmarshal(w.Body.Bytes(), &createdExercise)
					So(w.Code, ShouldEqual, http.StatusCreated)
					So(createdExercise.Name, ShouldEqual, exercise.Name)
					So(createdExercise.ID, ShouldNotBeEmpty)
				})
			})
			Convey("And a name is not given", func() {
				exercise.Name = ""
				body, _ := json.Marshal(exercise)
				req, _ := http.NewRequest("POST", "/exercises", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				Convey("Then an error is returned", func() {
					r.ServeHTTP(w, req)
					// Parse response
					var response map[string]string
					err := json.Unmarshal(w.Body.Bytes(), &response)
					So(err, ShouldBeNil)
					So(response["error"], ShouldEqual, "Name cannot be empty")
				})
			})
		})
	})
}

func TestGetExerciseByID(t *testing.T) {
	Convey("Given a database and an exercise", t, func() {
		db, cleanup := setupTestDB(t)
		defer cleanup()
		testExercise := models.Exercise{
			Name: "Bench Press",
		}
		db.Create(&testExercise)
		w := httptest.NewRecorder()
		r := setupRouter(db)
		Convey("When getting an exercise", func() {
			Convey("And the ID exists", func() {
				req, _ := http.NewRequest("GET", "/exercises/"+strconv.Itoa(int(testExercise.ID)), nil)
				req.Header.Set("Content-Type", "application/json")
				Convey("Then the exercise is returned", func() {
					r.ServeHTTP(w, req)
					var fetchedExercise models.Exercise
					json.Unmarshal(w.Body.Bytes(), &fetchedExercise)
					So(w.Code, ShouldEqual, http.StatusOK)
					So(fetchedExercise.Name, ShouldEqual, testExercise.Name)
					So(fetchedExercise.ID, ShouldEqual, testExercise.ID)
				})
			})
			Convey("And the ID does not exist", func() {
				req, _ := http.NewRequest("GET", "/exercises/10", nil)
				req.Header.Set("Content-Type", "application/json")
				Convey("Then an error is returned", func() {
					r.ServeHTTP(w, req)
					// Parse response
					var response map[string]string
					err := json.Unmarshal(w.Body.Bytes(), &response)
					So(err, ShouldBeNil)
					So(response["error"], ShouldEqual, "record not found")
				})
			})
		})
	})
}

func TestUpdateExercise(t *testing.T) {
	Convey("Given a database and an exercise", t, func() {
		db, cleanup := setupTestDB(t)
		defer cleanup()
		testExercise := models.Exercise{
			Name: "Bench Press",
		}
		db.Create(&testExercise)
		w := httptest.NewRecorder()
		r := setupRouter(db)
		Convey("When updating an exercise", func() {
			exercise := models.Exercise{
				Name: "Squat",
			}
			body, _ := json.Marshal(exercise)
			Convey("And the ID exists", func() {
				req, _ := http.NewRequest("PUT", "/exercises/"+strconv.Itoa(int(testExercise.ID)), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				Convey("Then the updated exercise is returned", func() {
					r.ServeHTTP(w, req)
					var updatedExercise models.Exercise
					json.Unmarshal(w.Body.Bytes(), &updatedExercise)
					So(w.Code, ShouldEqual, http.StatusOK)
					So(updatedExercise.Name, ShouldEqual, exercise.Name)
				})
			})
			Convey("And the ID does not exist", func() {
				req, _ := http.NewRequest("PUT", "/exercises/10", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				Convey("Then an error is returned", func() {
					r.ServeHTTP(w, req)
					// Parse response
					var response map[string]string
					err := json.Unmarshal(w.Body.Bytes(), &response)
					So(err, ShouldBeNil)
					So(response["error"], ShouldEqual, "record not found")
				})
			})
		})
	})
}

func TestDeleteExercise(t *testing.T) {
	Convey("Given a database and an exercise", t, func() {
		db, cleanup := setupTestDB(t)
		defer cleanup()
		testExercise := models.Exercise{
			Name: "Bench Press",
		}
		db.Create(&testExercise)
		w := httptest.NewRecorder()
		r := setupRouter(db)
		Convey("When deleting an exercise", func() {
			Convey("And the ID exists", func() {
				req, _ := http.NewRequest("DELETE", "/exercises/"+strconv.Itoa(int(testExercise.ID)), nil)
				req.Header.Set("Content-Type", "application/json")
				Convey("Then the exercise is deleted", func() {
					r.ServeHTTP(w, req)
					var response map[string]string
					err := json.Unmarshal(w.Body.Bytes(), &response)
					So(err, ShouldBeNil)
					So(response["message"], ShouldEqual, "Exercise deleted")
				})
			})
			Convey("And an invalid ID is provided", func() {
				req, _ := http.NewRequest("DELETE", "/exercises/invalid", nil)
				req.Header.Set("Content-Type", "application/json")
				Convey("Then an error is returned", func() {
					r.ServeHTTP(w, req)
					var response map[string]string
					err := json.Unmarshal(w.Body.Bytes(), &response)
					So(err, ShouldBeNil)
					So(response["error"], ShouldEqual, "invalid exercise ID")
				})
			})
		})
	})
}
