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

	ErrMissingToken = New(
		CodeAuthError,
		"missing token",
	)

	ErrInvalidToken = New(
		CodeAuthError,
		"invalid token",
	)
)
