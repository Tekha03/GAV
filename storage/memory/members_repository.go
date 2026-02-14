package memory

import (
	"context"
	"errors"
	"gav/internal/chat"
	"sync"
)

var (
	ErrMemberExists = errors.New("chat member already exists")
	ErrMemberNotFound = errors.New("member not found")
)

type MembersRepository struct {
	mu 		sync.RWMutex
	members map[uint]map[uint]*chat.ChatMember
}

func NewMembersRepository() *MembersRepository {
	return &MembersRepository{members: make(map[uint]map[uint]*chat.ChatMember)}
}

func (mr *MembersRepository) AddMember(ctx context.Context, member *chat.ChatMember) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[member.ChatID]; !ok {
		mr.members[member.ChatID] = make(map[uint]*chat.ChatMember)
	}

	if _, exists := mr.members[member.ChatID][member.UserID]; exists {
		return ErrMemberExists
	}

	mr.members[member.ChatID][member.UserID] = member
	return nil
}

func (mr *MembersRepository) RemoveMember(ctx context.Context, memberID, chatID uint) error {
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

func (mr *MembersRepository) GetMembers(ctx context.Context, chatID uint) ([]*chat.ChatMember, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	membersMap, ok := mr.members[chatID]
	if !ok {
		return []*chat.ChatMember{}, nil
	}

	result := make([]*chat.ChatMember, 0, len(membersMap))

	for _, member := range membersMap {
		result = append(result, member)
	}

	return result, nil
}

func (mr *MembersRepository) UpdateRole(ctx context.Context, chatID, userID uint, role *chat.MemberRole) error {
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

func (mr *MembersRepository) SetMuted(ctx context.Context, chatID, userID uint, muted bool) error {
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