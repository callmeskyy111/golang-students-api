package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/callmeskyy111/golang-students-api/internal/types"
	"github.com/callmeskyy111/golang-students-api/internal/utils/respnse"
)

func New() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

		var student types.Student

		// Decode json data
		err:=json.NewDecoder(r.Body).Decode(&student)
		// Check specific error
		if errors.Is(err, io.EOF){
			respnse.WriteJson(w, http.StatusBadRequest, respnse.GneralError(fmt.Errorf("empty body")))
			return 
		}

		slog.Info("Creating a student..")
		respnse.WriteJson(w, http.StatusCreated, map[string]string{"success":"Ok ✅"})
		//w.Write([]byte("Welcome to Students-API ✅"))
	}
}