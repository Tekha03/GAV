package gorm

import (
	"context"
	"errors"
	"messenger/internal/model"
	"messenger/internal/repository"
	apperrors "shared/app_errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatMemberRepository struct {
	repo *Repository
}

func NewChatMemberRepository(repo *Repository) repository.ChatMemberRepository {
	return &ChatMemberRepository{repo: repo}
}

func (cmr *ChatMemberRepository) AddMember(ctx context.Context, member *model.ChatMember) error {
	err := cmr.repo.WithContext(ctx).Create(member).Error
	return createError(err, apperrors.ChatMemberAlreadyExists, "chat member already exists", "failed to add chat member")
}

func (cmr *ChatMemberRepository) RemoveMember(ctx context.Context, memberID, chatID uuid.UUID) error {
	result := cmr.repo.WithContext(ctx).
		Where("chat_id = ? AND user_id = ?", chatID, memberID).
		Delete(&model.ChatMember{})
	return mutationError(result, apperrors.ChatMemberNotFound, "chat member not found", "failed to remove chat member")
}

func (cmr *ChatMemberRepository) GetMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error) {
	var members []*model.ChatMember
	err := cmr.repo.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Find(&members).Error
	if err != nil {
		return nil, internalError("failed to get chat members", err)
	}
	return members, nil
}

func (cmr *ChatMemberRepository) UpdateRole(ctx context.Context, chatID, userID uuid.UUID, role *model.MemberRole) error {
	result := cmr.repo.WithContext(ctx).
		Model(&model.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("role", *role)
	return mutationError(result, apperrors.ChatMemberNotFound, "chat member not found", "failed to update chat member role")
}

func (cmr *ChatMemberRepository) SetMuted(ctx context.Context, chatID, userID uuid.UUID, muted bool) error {
	result := cmr.repo.WithContext(ctx).
		Model(&model.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("muted", muted)
	return mutationError(result, apperrors.ChatMemberNotFound, "chat member not found", "failed to update chat member mute state")
}

func (cmr *ChatMemberRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var chats []struct {
		ChatID uuid.UUID `gorm:"column:chat_id"`
	}
	err := cmr.repo.WithContext(ctx).
		Table("chat_members").
		Where("user_id = ?", userID).
		Select("chat_id").
		Find(&chats).Error

	if err != nil {
		return nil, internalError("failed to get user chats", err)
	}

	result := make([]uuid.UUID, len(chats))
	for i, c := range chats {
		result[i] = c.ChatID
	}
	return result, nil
}

func (cmr *ChatMemberRepository) FindPrivateChatBetween(ctx context.Context, userID1, userID2 uuid.UUID) (uuid.UUID, error) {
	var row struct {
		ChatID uuid.UUID `gorm:"column:chat_id"`
	}

	err := cmr.repo.WithContext(ctx).
		Table("chat_members AS cm1").
		Select("cm1.chat_id").
		Joins("JOIN chat_members AS cm2 ON cm2.chat_id = cm1.chat_id").
		Joins("JOIN chats ON chats.id = cm1.chat_id").
		Where("chats.is_group = false").
		Where("cm1.user_id = ? AND cm2.user_id = ?", userID1, userID2).
		Limit(1).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, nil
		}
		return uuid.Nil, internalError("failed to find private chat", err)
	}

	return row.ChatID, nil
}

func (cmr *ChatMemberRepository) GetLastReadMessageID(ctx context.Context, chatID, userID uuid.UUID) (uuid.UUID, error) {
	var member model.ChatMember
	err := cmr.repo.WithContext(ctx).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		First(&member).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
		}
		return uuid.Nil, internalError("failed to get last read message ID", err)
	}

	return member.LastReadMessageID, nil
}

func (cmr *ChatMemberRepository) UpdateLastReadMessageID(ctx context.Context, chatID, userID, messageID uuid.UUID) error {
	result := cmr.repo.WithContext(ctx).
		Model(&model.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("last_read_message_id", messageID)

	return mutationError(result, apperrors.ChatMemberNotFound, "chat member not found", "failed to update last read message")
}

func (cmr *ChatMemberRepository) MemberExists(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (bool, error) {
	members, err := cmr.GetMembers(ctx, chatID)
	if err != nil {
		return false, err
	}

	ok := false
	for _, member := range members {
		if userID == member.UserID {
			ok = true
		}
	}

	return ok, nil
}

func (cmr *ChatMemberRepository) GetRole(ctx context.Context, userID, chatID uuid.UUID) (*model.MemberRole, error) {
	var member model.ChatMember
	err := cmr.repo.WithContext(ctx).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		First(&member).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
		}
		return nil, internalError("failed to get chat member role", err)
	}

	return &member.Role, nil
}
