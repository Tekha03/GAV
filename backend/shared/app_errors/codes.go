package apperrors

type Code string

// Common errors
const (
	CodeInternal           Code = "INTERNAL_ERROR"
	CodeValidation         Code = "VALIDATION_ERROR"
	CodeResourceNotFound   Code = "RESOURCE_NOT_FOUND"
	CodeServiceUnavailable Code = "SERVICE_UNAVAILABLE"
	CodeRequestCancelled   Code = "REQUEST_CANCELLED"
	CodeRequestTimeout     Code = "REQUEST_TIMEOUT"
)

// Auth and tokens
const (
	CodeAuthCredentialsInvalid Code = "AUTH_CREDENTIALS_INVALID"
	CodeAuthTokenMissing       Code = "AUTH_TOKEN_MISSING"
	CodeAuthTokenInvalid       Code = "AUTH_TOKEN_INVALID"
	CodeAuthTokenExpired       Code = "AUTH_TOKEN_EXPIRED"
	CodeAuthRefreshInvalid     Code = "AUTH_REFRESH_TOKEN_INVALID"
	CodeAuthForbidden          Code = "AUTH_FORBIDDEN"
	CodeAuthEmailExists        Code = "AUTH_EMAIL_ALREADY_EXISTS"
)

// Users and profiles
const (
	CodeUserNotFound           Code = "USER_NOT_FOUND"
	CodeUserEmailRequired      Code = "USER_EMAIL_REQUIRED"
	CodeUserPasswordRequired   Code = "USER_PASSWORD_REQUIRED"
	CodeProfileNotFound        Code = "PROFILE_NOT_FOUND"
	CodeProfileAlreadyExists   Code = "PROFILE_ALREADY_EXISTS"
	CodeProfileLocationInvalid Code = "PROFILE_LOCATION_INVALID"
)

// Dogs
const (
	CodeDogNotFound      Code = "DOG_NOT_FOUND"
	CodeDogAlreadyExists Code = "DOG_ALREADY_EXISTS"
	CodeDogAccessDenied  Code = "DOG_ACCESS_DENIED"
	CodeDogNameRequired  Code = "DOG_NAME_REQUIRED"
	CodeDogBreedRequired Code = "DOG_BREED_REQUIRED"
	CodeDogGenderInvalid Code = "DOG_GENDER_INVALID"
	CodeDogStatusInvalid Code = "DOG_STATUS_INVALID"
	CodeDogAgeInvalid    Code = "DOG_AGE_INVALID"
)

// Posts and interactions
const (
	CodePostNotFound        Code = "POST_NOT_FOUND"
	CodePostContentRequired Code = "POST_CONTENT_REQUIRED"
	CodePostAccessDenied    Code = "POST_ACCESS_DENIED"

	CodeCommentNotFound     Code = "COMMENT_NOT_FOUND"
	CodeCommentAccessDenied Code = "COMMENT_ACCESS_DENIED"

	CodeLikeAlreadyExists Code = "LIKE_ALREADY_EXISTS"
	CodeLikeNotFound      Code = "LIKE_NOT_FOUND"

	CodeFollowSelfNotAllowed Code = "FOLLOW_SELF_NOT_ALLOWED"
	CodeFollowAlreadyExists  Code = "FOLLOW_ALREADY_EXISTS"
	CodeFollowNotFound       Code = "FOLLOW_NOT_FOUND"
)

// Messenger
const (
	CodeChatNotFound       Code = "CHAT_NOT_FOUND"
	CodeChatAlreadyExists  Code = "CHAT_ALREADY_EXISTS"
	CodeChatAccessDenied   Code = "CHAT_ACCESS_DENIED"
	CodeChatSelfNotAllowed Code = "CHAT_SELF_NOT_ALLOWED"
	CodeChatGroupRequired  Code = "CHAT_GROUP_REQUIRED"

	CodeChatMemberNotFound      Code = "CHAT_MEMBER_NOT_FOUND"
	CodeChatMemberAlreadyExists Code = "CHAT_MEMBER_ALREADY_EXISTS"

	CodeMessageNotFound                 Code = "MESSAGE_NOT_FOUND"
	CodeMessageAlreadyExists            Code = "MESSAGE_ALREADY_EXISTS"
	CodeMessageContentRequired          Code = "MESSAGE_CONTENT_REQUIRED"
	CodeMessageTextTooLong              Code = "MESSAGE_TEXT_TOO_LONG"
	CodeMessageAttachmentsLimitExceeded Code = "MESSAGE_ATTACHMENTS_LIMIT_EXCEEDED"
	CodeMessageReplyInvalid             Code = "MESSAGE_REPLY_INVALID"

	CodeAttachmentNotFound      Code = "ATTACHMENT_NOT_FOUND"
	CodeAttachmentAlreadyExists Code = "ATTACHMENT_ALREADY_EXISTS"
	CodeReactionNotFound        Code = "REACTION_NOT_FOUND"
	CodeReactionAlreadyExists   Code = "REACTION_ALREADY_EXISTS"
)
