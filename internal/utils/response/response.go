package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"` //struct tags
	Error  string `json:"error"`
}

const (
	StatusSuccess = "success"
	StatusFailed = "failed"
)

func WriteJson(w http.ResponseWriter, statusCode int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

func GeneralErr(err error) Response {
	return Response{
		Status: StatusFailed,
		Error: err.Error(),
	}
}

func ValidationErr(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err:= range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response {
		Status: StatusFailed,
		Error: strings.Join(errMsgs, ", "),
	}

}