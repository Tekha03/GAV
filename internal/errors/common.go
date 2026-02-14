package errors

type Code string

const (
	CodeUnauthorized	Code = "UNAUTHORIZED"
	CodeForbidden		Code = "FORBIDDEN"
	CodeNotFound		Code = "NOT_FOUND"
	CodeConflict		 Code = "CONFLICT"
	CodeValidation		Code = "VALIDATION_ERROR"
	CodeInternal		Code = "INTERNAL_ERROR"
	CodeAuthError		Code = "AUTH_ERROR"
)

type Error struct {
	Code	Code
	Message	string
}

func New(code Code, message string) *Error {
	return &Error{
		Code: code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
