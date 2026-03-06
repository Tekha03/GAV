package sqlite

import (
	"context"
	"errors"
	"gav/internal/chat/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembersRepository struct {
	*BaseRepository
}

func NewMembersRepository(db *gorm.DB) (*MembersRepository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &MembersRepository{BaseRepository: repo}, nil
}

func (r *MembersRepository) AddMember(ctx context.Context, member *model.ChatMember) error {
	return r.DB(ctx).Create(member).Error
}

func (r *MembersRepository) RemoveMember(ctx context.Context, chatID, userID uuid.UUID) error {
	res := r.DB(ctx).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Delete(&model.ChatMember{})

	if res.RowsAffected == 0 {
		return ErrMemberNotFound
	}
	return res.Error
}

func (r *MembersRepository) GetMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error) {
	var members []*model.ChatMember
	err := r.DB(ctx).
		Where("chat_id = ?", chatID).
		Find(&members).Error

	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *MembersRepository) UpdateRole(ctx context.Context, chatID, userID uuid.UUID, role *model.MemberRole) error {
	res := r.DB(ctx).
		Model(&model.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("role", *role)

	if res.RowsAffected == 0 {
		return ErrMemberNotFound
	}
	return res.Error
}

func (r *MembersRepository) SetMuted(ctx context.Context, chatID, userID uuid.UUID, muted bool) error {
	res := r.DB(ctx).
		Model(&model.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("muted", muted)

	if res.RowsAffected == 0 {
		return ErrMemberNotFound
	}
	return res.Error
}

func (r *MembersRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var members []*model.ChatMember
	err := r.DB(ctx).
		Where("user_id = ?", userID).
		Find(&members).Error
	if err != nil {
		return nil, err
	}

	var chatIDs []uuid.UUID
	for _, m := range members {
		chatIDs = append(chatIDs, m.ChatID)
	}

	return chatIDs, nil
}

func (r *MembersRepository) GetLastReadMessageID(ctx context.Context, chatID, userID uuid.UUID) (uuid.UUID, error) {
	var member model.ChatMember
	err := r.DB(ctx).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		First(&member).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, ErrMemberNotFound
	}
	if err != nil {
		return uuid.Nil, err
	}

	return member.LastReadMessageID, nil
}