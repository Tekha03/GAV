package memory

import (
	"context"
	"messenger/internal/model"
	"messenger/internal/repository"
	"sync"

	"github.com/google/uuid"
)

type MembersRepository struct {
	mu 		sync.RWMutex
	members map[uuid.UUID]map[uuid.UUID]*model.ChatMember
}

func NewMembersRepository() repository.ChatMemberRepository {
	return &MembersRepository{members: make(map[uuid.UUID]map[uuid.UUID]*model.ChatMember)}
}

func (mr *MembersRepository) AddMember(ctx context.Context, member *model.ChatMember) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[member.ChatID]; !ok {
		mr.members[member.ChatID] = make(map[uuid.UUID]*model.ChatMember)
	}

	if _, exists := mr.members[member.ChatID][member.UserID]; exists {
		return repository.ErrMemberExists
	}

	mr.members[member.ChatID][member.UserID] = member
	return nil
}

func (mr *MembersRepository) RemoveMember(ctx context.Context, memberID, chatID uuid.UUID) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[chatID]; !ok {
		return repository.ErrMemberNotFound
	}

	if _, exists := mr.members[chatID][memberID]; !exists {
		return repository.ErrMemberNotFound
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
		return repository.ErrMemberNotFound
	}

	if _, exists := mr.members[chatID][userID]; !exists {
		return repository.ErrMemberNotFound
	}

	mr.members[chatID][userID].Role = *role
	return nil
}

func (mr *MembersRepository) SetMuted(ctx context.Context, chatID, userID uuid.UUID, muted bool) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, ok := mr.members[chatID]; !ok {
		return repository.ErrMemberNotFound
	}

	if _, exists := mr.members[chatID][userID]; !exists {
		return repository.ErrMemberNotFound
	}

	mr.members[chatID][userID].Muted = muted
	return nil
}

func (mr *MembersRepository) GetLastReadMessageID(ctx context.Context, chatID, userID uuid.UUID) (uuid.UUID, error) {
    mr.mu.RLock()
    defer mr.mu.RUnlock()

    membersMap, ok := mr.members[chatID]
    if !ok {
        return uuid.Nil, repository.ErrChatNotFound
    }

    member, exists := membersMap[userID]
    if !exists {
        return uuid.Nil, repository.ErrMemberNotFound
    }

    return member.LastReadMessageID, nil
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
