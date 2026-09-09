package memory

import (
	"context"
	"messenger/internal/model"
	"messenger/internal/repository"
	apperrors "shared/app_errors"
	"sync"

	"github.com/google/uuid"
)

type MembersRepository struct {
	mu      sync.RWMutex
	members map[uuid.UUID]map[uuid.UUID]*model.ChatMember
}

func NeMembersRepository() repository.ChatMemberRepository {
	return &MembersRepository{members: make(map[uuid.UUID]map[uuid.UUID]*model.ChatMember)}
}

func (mr *MembersRepository) AddMember(ctx context.Context, member *model.ChatMember) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[member.ChatID]; !ok {
		mr.members[member.ChatID] = make(map[uuid.UUID]*model.ChatMember)
	}

	if _, exists := mr.members[member.ChatID][member.UserID]; exists {
		return apperrors.New(apperrors.ChatMemberAlreadyExists, "chat member already exists")
	}

	mr.members[member.ChatID][member.UserID] = member
	return nil
}

func (mr *MembersRepository) RemoveMember(ctx context.Context, memberID, chatID uuid.UUID) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[chatID]; !ok {
		return apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	if _, exists := mr.members[chatID][memberID]; !exists {
		return apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	delete(mr.members[chatID], memberID)
	return nil
}

func (mr *MembersRepository) GetMembers(ctx context.Context, chatID uuid.UUID) ([]*model.ChatMember, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	membersMap, ok := mr.members[chatID]
	if !ok {
		return []*model.ChatMember{}, nil
	}

	result := make([]*model.ChatMember, 0, len(membersMap))

	for _, member := range membersMap {
		result = append(result, member)
	}

	return result, nil
}

func (mr *MembersRepository) UpdateRole(ctx context.Context, chatID, userID uuid.UUID, role *model.MemberRole) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[chatID]; !ok {
		return apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	if _, exists := mr.members[chatID][userID]; !exists {
		return apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	mr.members[chatID][userID].Role = *role
	return nil
}

func (mr *MembersRepository) SetMuted(ctx context.Context, chatID, userID uuid.UUID, muted bool) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[chatID]; !ok {
		return apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	if _, exists := mr.members[chatID][userID]; !exists {
		return apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	mr.members[chatID][userID].Muted = muted
	return nil
}

func (mr *MembersRepository) GetLastReadMessageID(ctx context.Context, chatID, userID uuid.UUID) (uuid.UUID, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	membersMap, ok := mr.members[chatID]
	if !ok {
		return uuid.Nil, apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	member, exists := membersMap[userID]
	if !exists {
		return uuid.Nil, apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	return member.LastReadMessageID, nil
}

func (mr *MembersRepository) UpdateLastReadMessageID(ctx context.Context, chatID, userID, messageID uuid.UUID) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	membersMap, ok := mr.members[chatID]
	if !ok {
		return apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	member, exists := membersMap[userID]
	if !exists {
		return apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	member.LastReadMessageID = messageID
	return nil
}

func (mr *MembersRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var chats []uuid.UUID
	for chatID, members := range mr.members {
		if _, exists := members[userID]; exists {
			chats = append(chats, chatID)
		}
	}
	return chats, nil
}

func (mr *MembersRepository) FindPrivateChatBetween(ctx context.Context, userID1, userID2 uuid.UUID) (uuid.UUID, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	for chatID, members := range mr.members {
		_, hasFirst := members[userID1]
		_, hasSecond := members[userID2]
		if hasFirst && hasSecond && len(members) == 2 {
			return chatID, nil
		}
	}

	return uuid.Nil, nil
}

func (mr *MembersRepository) MemberExists(ctx context.Context, userID, chatID uuid.UUID) (bool, error) {
	members, err := mr.GetMembers(ctx, chatID)
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

func (mr *MembersRepository) GetRole(ctx context.Context, userID, chatID uuid.UUID) (*model.MemberRole, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	membersMap, ok := mr.members[chatID]
	if !ok {
		return nil, apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	member, exists := membersMap[userID]
	if !exists {
		return nil, apperrors.New(apperrors.ChatMemberNotFound, "chat member not found")
	}

	return &member.Role, nil
}
