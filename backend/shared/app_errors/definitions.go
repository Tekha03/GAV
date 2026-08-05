package apperrors

// Common errors
var (
	Internal           = define(CodeInternal, CategoryInternal)
	Validation         = define(CodeValidation, CategoryValidation)
	ResourceNotFound   = define(CodeResourceNotFound, CategoryNotFound)
	ServiceUnavailable = define(CodeServiceUnavailable, CategoryUnavailable)
	RequestCancelled   = define(CodeRequestCancelled, CategoryCancelled)
	RequestTimeout     = define(CodeRequestTimeout, CategoryUnavailable)
)

// Auth and tokens
var (
	AuthCredentialsInvalid = define(CodeAuthCredentialsInvalid, CategoryUnauthenticated)
	AuthTokenMissing       = define(CodeAuthTokenMissing, CategoryUnauthenticated)
	AuthTokenInvalid       = define(CodeAuthTokenInvalid, CategoryUnauthenticated)
	AuthTokenExpired       = define(CodeAuthTokenExpired, CategoryUnauthenticated)
	AuthRefreshInvalid     = define(CodeAuthRefreshInvalid, CategoryUnauthenticated)
	AuthForbidden          = define(CodeAuthForbidden, CategoryPermissionDenied)
	AuthEmailExists        = define(CodeAuthEmailExists, CategoryConflict)
)

// Users and profiles
var (
	UserNotFound           = define(CodeUserNotFound, CategoryNotFound)
	UserEmailRequired      = define(CodeUserEmailRequired, CategoryValidation)
	UserPasswordRequired   = define(CodeUserPasswordRequired, CategoryValidation)
	ProfileNotFound        = define(CodeProfileNotFound, CategoryNotFound)
	ProfileAlreadyExists   = define(CodeProfileAlreadyExists, CategoryConflict)
	ProfileLocationInvalid = define(CodeProfileLocationInvalid, CategoryValidation)
)

// Dogs
var (
	DogNotFound      = define(CodeDogNotFound, CategoryNotFound)
	DogAlreadyExists = define(CodeDogAlreadyExists, CategoryConflict)
	DogAccessDenied  = define(CodeDogAccessDenied, CategoryPermissionDenied)
	DogNameRequired  = define(CodeDogNameRequired, CategoryValidation)
	DogBreedRequired = define(CodeDogBreedRequired, CategoryValidation)
	DogGenderInvalid = define(CodeDogGenderInvalid, CategoryValidation)
	DogStatusInvalid = define(CodeDogStatusInvalid, CategoryValidation)
	DogAgeInvalid    = define(CodeDogAgeInvalid, CategoryValidation)
)

// Posts and interactions
var (
	PostNotFound        = define(CodePostNotFound, CategoryNotFound)
	PostContentRequired = define(CodePostContentRequired, CategoryValidation)
	PostAccessDenied    = define(CodePostAccessDenied, CategoryPermissionDenied)

	CommentNotFound     = define(CodeCommentNotFound, CategoryNotFound)
	CommentAccessDenied = define(CodeCommentAccessDenied, CategoryPermissionDenied)

	LikeAlreadyExists = define(CodeLikeAlreadyExists, CategoryConflict)
	LikeNotFound      = define(CodeLikeNotFound, CategoryNotFound)

	FollowSelfNotAllowed = define(CodeFollowSelfNotAllowed, CategoryValidation)
	FollowAlreadyExists  = define(CodeFollowAlreadyExists, CategoryConflict)
	FollowNotFound       = define(CodeFollowNotFound, CategoryNotFound)
)

// Messenger
var (
	ChatNotFound       = define(CodeChatNotFound, CategoryNotFound)
	ChatAlreadyExists  = define(CodeChatAlreadyExists, CategoryConflict)
	ChatAccessDenied   = define(CodeChatAccessDenied, CategoryPermissionDenied)
	ChatSelfNotAllowed = define(CodeChatSelfNotAllowed, CategoryValidation)

	ChatMemberNotFound      = define(CodeChatMemberNotFound, CategoryNotFound)
	ChatMemberAlreadyExists = define(CodeChatMemberAlreadyExists, CategoryConflict)

	MessageNotFound        = define(CodeMessageNotFound, CategoryNotFound)
	MessageContentRequired = define(CodeMessageContentRequired, CategoryValidation)
	MessageTextTooLong     = define(CodeMessageTextTooLong, CategoryValidation)

	AttachmentNotFound = define(CodeAttachmentNotFound, CategoryNotFound)
	ReactionNotFound   = define(CodeReactionNotFound, CategoryNotFound)
)
