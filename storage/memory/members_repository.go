package memory

import (
	"context"
	"errors"
	"gav/internal/chat/model"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrMemberExists = errors.New("chat member already exists")
	ErrMemberNotFound = errors.New("member not found")
)

type MembersRepository struct {
	mu 		sync.RWMutex
	members map[uuid.UUID]map[uuid.UUID]*model.ChatMember
}

func NewMembersRepository() *MembersRepository {
	return &MembersRepository{members: make(map[uuid.UUID]map[uuid.UUID]*model.ChatMember)}
}

func (mr *MembersRepository) AddMember(ctx context.Context, member *model.ChatMember) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[member.ChatID]; !ok {
		mr.members[member.ChatID] = make(map[uuid.UUID]*model.ChatMember)
	}

	if _, exists := mr.members[member.ChatID][member.UserID]; exists {
		return ErrMemberExists
	}

	mr.members[member.ChatID][member.UserID] = member
	return nil
}

func (mr *MembersRepository) RemoveMember(ctx context.Context, memberID, chatID uuid.UUID) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[chatID]; !ok {
		return ErrMemberNotFound
	}

	if _, exists := mr.members[chatID][memberID]; !exists {
		return ErrMemberNotFound
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
		return ErrMemberNotFound
	}

	if _, exists := mr.members[chatID][userID]; !exists {
		return ErrMemberNotFound
	}

	mr.members[chatID][userID].Role = *role
	return nil
}

func (mr *MembersRepository) SetMuted(ctx context.Context, chatID, userID uuid.UUID, muted bool) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[chatID]; !ok {
		return ErrMemberNotFound
	}

	if _, exists := mr.members[chatID][userID]; !exists {
		return ErrMemberNotFound
	}

	mr.members[chatID][userID].Muted = muted
	return nil
}

func (mr *MembersRepository) GetUserChats(ctx context.Context, userId uuid.UUID) ([]*uuid.UUID, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var result []*uuid.UUID
	for chatID, chats := range mr.members {
		for id := range chats {
			if id == userId {
				result = append(result, &chatID)
			}
		}
	}

	return result, nil
}