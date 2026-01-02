package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/student-ankitpandit/rest-api/http/handlers/student"
	"github.com/student-ankitpandit/rest-api/internal/storage"
	"github.com/student-ankitpandit/rest-api/internal/types"
	"github.com/student-ankitpandit/rest-api/internal/utils/response"
)



func New(storage storage.Storage) http.HandlerFunc { //this is called dependency injection
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
		
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
	
		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if(err != nil) {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetStudentsById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a user through their id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("something went wrong", slog.String("id", id))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralErr(err))
			return 
		}
		student, err := storage.GetStudentsById(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralErr(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		students, err := storage.GetStudentsList()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

