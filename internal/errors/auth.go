package errors

var (
	ErrInvalidCredentials = New(
		CodeUnauthorized,
		"invalid email or password",
	)

	ErrAlreadyExists = New(
		CodeConflict,
		"user already exists",
	)
)
