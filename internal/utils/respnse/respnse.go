package respnse

import (
	"encoding/json"
	"net/http"
)

type Response struct{
	Status string `json:"status"`
	Error string `json:"error"`
}

const (
	StatusOk = "OK✅"
	StatusError = "Error 🔴"
)

func WriteJson(w http.ResponseWriter, status int, data any)error {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GneralError(err error)Response{
return Response{
	Status: StatusError,
	Error: err.Error(),
}
}