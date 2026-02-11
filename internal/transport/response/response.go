package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(writer http.ResponseWriter, code int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(writer).Encode(data)
	}
}

func Error(writer http.ResponseWriter, code int, err error) {
	JSON(writer, code, ErrorResponse{Error: err.Error()})
}
