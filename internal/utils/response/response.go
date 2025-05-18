package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOk                  = "ok"
	StatusCreated             = "created"
	StatusBadRequest          = "bad request"
	StatusInternalServerError = "internal server error"
	StatusError               = "error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status:  StatusError,
		Message: err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errs {

		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s': %s", err.Field(), err.Tag()))
		}
	}
	return Response{
		Status:  StatusBadRequest,
		Message: strings.Join(errMsgs, ", "),
	}
}
