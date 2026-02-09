package response

import "net/http"

var errorMap = map[string]mappedError{
	"UNAUTHORIZED": {
		status:  http.StatusUnauthorized,
		code:	 "UNAUTHORIZED",
		message: "unauthorized",
	},
	"FORBIDDEN": {
		status:  http.StatusForbidden,
		code: 	 "FORBIDDEN",
		message: "forbidden",
	},
	"NOT_FOUND": {
		status: http.StatusNotFound,
		code: 	"NOT_FOUND",
		message: "resource not found",
	},
	"CONFLICT": {
		status:  http.StatusConflict,
		code:	 "CONFLICT",
		message: "conflict",
	},
	"VALIDATION_ERROR": {
		status:  http.StatusBadRequest,
		code: 	 "VALIDATION_ERROR",
		message: "invalid request",
	},
}
