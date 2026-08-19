package sqlite

import apperrors "shared/app_errors"

var (
	ErrCommentNotFound      = apperrors.New(apperrors.CommentNotFound, "comment not found")
	ErrDogNotFound          = apperrors.New(apperrors.DogNotFound, "dog not found")
	ErrSettingsNotFound     = apperrors.New(apperrors.SettingsNotFound, "settings not found")
	ErrPostNotFound         = apperrors.New(apperrors.PostNotFound, "post not found")
	ErrVaccinationNotFound  = apperrors.New(apperrors.VaccinationNotFound, "vaccination not found")
	ErrUserNotFound         = apperrors.New(apperrors.UserNotFound, "user not found")
	ErrUserExists           = apperrors.New(apperrors.AuthEmailExists, "user already exists")
	ErrDBNil                = apperrors.New(apperrors.Internal, "database is nil")
	ErrRefreshTokenNotFound = apperrors.New(apperrors.AuthRefreshInvalid, "refresh token not found")

	ErrMemberExists   = apperrors.New(apperrors.ChatMemberAlreadyExists, "chat member already exists")
	ErrMemberNotFound = apperrors.New(apperrors.ChatMemberNotFound, "member not found")

	ErrChatExists   = apperrors.New(apperrors.ChatAlreadyExists, "chat already exists")
	ErrChatNotFound = apperrors.New(apperrors.ChatNotFound, "chat not found")
	ErrNotGroup     = apperrors.New(apperrors.ChatGroupRequired, "operation allowed only for group chats")
	ErrEmptyTitle   = apperrors.New(apperrors.Validation, "title cannot be empty")

	ErrMessageNotFound = apperrors.New(apperrors.MessageNotFound, "message not found")
	ErrMessageExists   = apperrors.New(apperrors.MessageAlreadyExists, "message already exists")

	ErrAttachmentNotFound = apperrors.New(apperrors.AttachmentNotFound, "attachment not found")
	ErrAttachmentExist    = apperrors.New(apperrors.AttachmentAlreadyExists, "attachment already exists")

	ErrReactionExists   = apperrors.New(apperrors.ReactionAlreadyExists, "reaction already exists")
	ErrReactionNotFound = apperrors.New(apperrors.ReactionNotFound, "reaction not found")
)
