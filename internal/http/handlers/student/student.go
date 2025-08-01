package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/callmeskyy111/golang-students-api/internal/storage"
	"github.com/callmeskyy111/golang-students-api/internal/types"
	"github.com/callmeskyy111/golang-students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

		slog.Info("Creating a student..")

		var student types.Student

		// Decode json data
		err:=json.NewDecoder(r.Body).Decode(&student)
		// Check specific error
		if errors.Is(err, io.EOF){
			response.WriteJson(w, http.StatusBadRequest, response.GneralError(fmt.Errorf("empty body")))
			return 
		}

		if err!= nil{
			response.WriteJson(w, http.StatusBadRequest, response.GneralError(err))
			return 
		}

		// request validation - golang validator package (Don't do it manually)
		if err:= validator.New().Struct(student); err!=nil{
			validateErrs:=err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest,response.ValidationErr(validateErrs))
			return 
		}

		lastId,err:=storage.CreateStudent(student.Name, student.Email, student.Age)

		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err!=nil{
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id":lastId})
		//w.Write([]byte("Welcome to Students-API ✅"))
	}
}

func GetById(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student.. ", slog.String("id",id))

		intId, err:= strconv.ParseInt(id,10,64)
		if err!=nil{
			slog.Error("error parsing the id")
			response.WriteJson(w, http.StatusBadRequest, response.GneralError(err))
			return 
		}
		student,err:= storage.GetStudentById(intId)
		if err!=nil{
			slog.Error("error getting user",slog.String("id",id))
			response.WriteJson(w, http.StatusInternalServerError, response.GneralError(err))
			return 
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students..")
		students, err:= storage.GetStudents()
		if err!= nil{
			slog.Error("Failed to get the student-list")
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		response.WriteJson(w, http.StatusOK, students)

	}
}
