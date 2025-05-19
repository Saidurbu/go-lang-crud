package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Saidurbu/go-lang-crud/internal/storage"
	"github.com/Saidurbu/go-lang-crud/internal/storage/sqlite"
	"github.com/Saidurbu/go-lang-crud/internal/types"
	"github.com/Saidurbu/go-lang-crud/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key") // store this in env in production

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type contextKey string

const emailContextKey = contextKey("email")

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty request body")))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validatorErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validatorErrs))
			return

		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Password, student.Age)

		slog.Info("Student created", "ID", slog.Int64("id", lastId))
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJSON(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "student created",
			"id":      lastId,
		})

	}

}

func Registration(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty request body")))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validatorErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validatorErrs))
			return

		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Password, student.Age)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "User registered"})
		slog.Info("User registered", "ID", slog.Int64("id", lastId))
	}

}

func Login(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&input)

		user, err := storage.GetStudentByEmail(input.Email)
		if err != nil {
			http.Error(w, "Invalid email", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		// Generate token
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Email: user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

	}
}

func Logout(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		int64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
			return
		}
		err = storage.DeleteStudent(int64)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]string{"success": "student deleted"})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		int64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
			return
		}
		student, err := storage.GetStudentById(int64)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, students)
	}
}
func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		int64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
			return
		}

		var student types.Student

		err = json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty request body")))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validatorErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validatorErrs))
			return

		}

		err = storage.UpdateStudent(int64, student.Name, student.Email, student.Password, student.Age)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]string{"success": "student updated"})
	}
}
func Delete(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		int64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
			return
		}
		err = storage.DeleteStudent(int64)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]string{"success": "student deleted"})
	}
}

func GetProfile(storage *sqlite.Sqlite) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		emailVal := ctx.Value(emailContextKey)
		email, ok := emailVal.(string)
		if !ok || email == "" {
			http.Error(w, "Email not found in token", http.StatusUnauthorized)
			return
		}

		student, err := storage.GetStudentByEmail(email)
		if err != nil {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		response := types.StudentResponse{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
			Age:   student.Age,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func EmailContextKey() interface{} {
	return emailContextKey
}
