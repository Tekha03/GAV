package apperrors

type Category string

const (
	CategoryValidation       Category = "validation"
	CategoryUnauthenticated  Category = "unauthenticated"
	CategoryPermissionDenied Category = "permission_denied"
	CategoryNotFound         Category = "not_found"
	CategoryConflict         Category = "conflict"
	CategoryRateLimited      Category = "rate_limited"
	CategoryUnavailable      Category = "unavailable"
	CategoryUnsupported      Category = "unsupported"
	CategoryCancelled        Category = "cancelled"
	CategoryInternal         Category = "internal"
)
