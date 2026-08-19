package response

import (
	"encoding/json"
	"net/http"
	"os"
	apperrors "shared/app_errors"

	"log/slog"
)

var logg = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

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

	appErr := apperrors.Normalize(err)
	message := appErr.Message
	if appErr.Category == apperrors.CategoryInternal {
		logg.Error("handler error", "error", err.Error())
		message = "internal server error"
	}

	JSON(w, httpStatus(appErr.Category), ErrorResponse{
		Error: ErrorBody{
			Code:     appErr.Code,
			Category: appErr.Category,
			Message:  message,
			Details:  appErr.Details,
		},
	})
}

func httpStatus(category apperrors.Category) int {
	switch category {
	case apperrors.CategoryValidation:
		return http.StatusBadRequest
	case apperrors.CategoryUnauthenticated:
		return http.StatusUnauthorized
	case apperrors.CategoryPermissionDenied:
		return http.StatusForbidden
	case apperrors.CategoryNotFound:
		return http.StatusNotFound
	case apperrors.CategoryConflict:
		return http.StatusConflict
	case apperrors.CategoryUnavailable:
		return http.StatusServiceUnavailable
	case apperrors.CategoryUnsupported:
		return http.StatusNotImplemented
	case apperrors.CategoryCancelled:
		return 499
	default:
		return http.StatusInternalServerError
	}
}

func InternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": "internal server error",
	})
}
