package response

import (
	"encoding/json"
	"net/http"

	"gav/internal/errors"
)

func JSON(writer http.ResponseWriter, code int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(writer).Encode(data)
	}
}

func Error(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if e, ok := err.(*errors.Error); ok {
		if mapped, exists := errorMap[string(e.Code)]; exists {
			JSON(w, mapped.status, ErrorResponse{
				Error: ErrorBody{
					Code: 	 mapped.code,
					Message: e.Message,
				},
			})
			return
		}
	}

	JSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: ErrorBody{
			Code: 	 "INTERNAL_ERROR",
			Message: "internal server error",
		},
	})
}
