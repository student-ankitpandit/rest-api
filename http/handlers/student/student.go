package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/student-ankitpandit/rest-api/internal/types"
	"github.com/student-ankitpandit/rest-api/internal/utils/response"
)



func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")
		
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			err := response.WriteJson(w, http.StatusBadRequest, response.GeneralErr(fmt.Errorf("empty body")))
			if(err != nil) {
				slog.Error("failed to write response", "error", err.Error())
			}
			return 
		}

		if(err != nil) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralErr(err))
			return
		}

		//req validation
		if err := validator.New().Struct(student); (err != nil) {

			// type assertion or type casting
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationErr(validateErrs))
			return 
		}
		


		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "true"})
	}

	
}

