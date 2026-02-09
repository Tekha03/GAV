package errors

var (
	ErrPostNotFound = New(
		CodeNotFound,
		"post not found",
	)

	ErrPostForbidden = New(
		CodeForbidden,
		"no permission to modify post",
	)
)
