package response

import (
	"encoding/json"
	"net/http"

	"social_network/internal/errors"

	"log/slog"
)

func JSON(writer http.ResponseWriter, code int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(writer).Encode(data)
	}
}

var logg slog.Logger

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

	logg.Error("handler error", "error", err.Error())
	JSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: ErrorBody{
			Code: 	 "INTERNAL_ERROR",
			Message: "internal server error",
		},
	})
}

func InternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": "internal server error",
	})
}
