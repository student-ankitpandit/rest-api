package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/student-ankitpandit/rest-api/internal/types"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		errors.Is(err, io.EOF)

		slog.Info("creating a student")

		w.Write([]byte("welcome to students api"))
	}

	
}