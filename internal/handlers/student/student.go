package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Saidurbu/go-lang-crud/internal/storage"
	"github.com/Saidurbu/go-lang-crud/internal/types"
	"github.com/Saidurbu/go-lang-crud/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

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

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		slog.Info("Student created", "ID", slog.Int64("id", lastId))
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastId})
		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "student created"})
	}

}
