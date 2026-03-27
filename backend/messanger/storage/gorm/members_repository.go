package gorm

import (
	"context"
	"errors"
	"messanger/internal/model"
	"messanger/internal/repository"

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
    return cmr.repo.WithContext(ctx).Create(member).Error
}

func (cmr *ChatMemberRepository) RemoveMember(ctx context.Context, memberID, chatID uuid.UUID) error {
    return cmr.repo.WithContext(ctx).
        Where("chat_id = ? AND user_id = ?", chatID, memberID).
        Delete(&model.ChatMember{}).Error
}

func (cmr *ChatMemberRepository) GetMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error) {
    var members []*model.ChatMember
    return members, cmr.repo.WithContext(ctx).
        Where("chat_id = ?", chatID).
        Find(&members).Error
}

func (cmr *ChatMemberRepository) UpdateRole(ctx context.Context, chatID, userID uuid.UUID, role *model.MemberRole) error {
    return cmr.repo.WithContext(ctx).
        Model(&model.ChatMember{}).
        Where("chat_id = ? AND user_id = ?", chatID, userID).
        Update("role", *role).Error
}

func (cmr *ChatMemberRepository) SetMuted(ctx context.Context, chatID, userID uuid.UUID, muted bool) error {
    return cmr.repo.WithContext(ctx).
        Model(&model.ChatMember{}).
        Where("chat_id = ? AND user_id = ?", chatID, userID).
        Update("muted", muted).Error
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
        return nil, err
    }
    
    result := make([]uuid.UUID, len(chats))
    for i, c := range chats {
        result[i] = c.ChatID
    }
    return result, nil
}

func (cmr *ChatMemberRepository) GetLastReadMessageID(ctx context.Context, chatID, userID uuid.UUID) (uuid.UUID, error) {
    var member model.ChatMember
    err := cmr.repo.WithContext(ctx).
        Where("chat_id = ? AND user_id = ?", chatID, userID).
        First(&member).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return uuid.Nil, nil
        }
        return uuid.Nil, err
    }
    
    return member.LastReadMessageID, nil
}
