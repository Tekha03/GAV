package service

import (
	"context"
	"sync"
	"testing"

	"messenger/internal/model"
	"messenger/internal/repository"

	"github.com/google/uuid"
)

type serializedTransactionManager struct {
	mu sync.Mutex
}

func (m *serializedTransactionManager) WithPrivateChatLock(ctx context.Context, _ string, fn func(context.Context) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return fn(ctx)
}

type privateChatRepository struct {
	repository.ChatRepository
	chats map[uuid.UUID]*model.Chat
}

func (r *privateChatRepository) Create(_ context.Context, chat *model.Chat) error {
	r.chats[chat.ID] = chat
	return nil
}

func (r *privateChatRepository) GetByID(_ context.Context, chatID uuid.UUID) (*model.Chat, error) {
	return r.chats[chatID], nil
}

func (r *privateChatRepository) GetByPrivateKey(_ context.Context, privateKey string) (*model.Chat, error) {
	for _, chat := range r.chats {
		if chat.PrivateKey != nil && *chat.PrivateKey == privateKey {
			return chat, nil
		}
	}
	return nil, nil
}

type privateMembersRepository struct {
	repository.ChatMemberRepository
	members map[uuid.UUID]map[uuid.UUID]*model.ChatMember
}

func (r *privateMembersRepository) AddMember(_ context.Context, member *model.ChatMember) error {
	if r.members[member.ChatID] == nil {
		r.members[member.ChatID] = make(map[uuid.UUID]*model.ChatMember)
	}
	r.members[member.ChatID][member.UserID] = member
	return nil
}

func (r *privateMembersRepository) FindPrivateChatBetween(_ context.Context, userID1, userID2 uuid.UUID) (uuid.UUID, error) {
	for chatID, members := range r.members {
		if members[userID1] != nil && members[userID2] != nil && len(members) == 2 {
			return chatID, nil
		}
	}
	return uuid.Nil, nil
}

func TestCreatePrivateChatIsIdempotentForConcurrentAndReversedRequests(t *testing.T) {
	chatRepo := &privateChatRepository{chats: make(map[uuid.UUID]*model.Chat)}
	membersRepo := &privateMembersRepository{members: make(map[uuid.UUID]map[uuid.UUID]*model.ChatMember)}
	svc := &ChatService{
		transactionManager: &serializedTransactionManager{},
		chatRepo:           chatRepo,
		membersRepo:        membersRepo,
	}
	userID1 := uuid.New()
	userID2 := uuid.New()

	const requests = 20
	results := make(chan uuid.UUID, requests)
	errors := make(chan error, requests)
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(reverse bool) {
			defer wg.Done()
			first, second := userID1, userID2
			if reverse {
				first, second = second, first
			}
			chat, err := svc.CreatePrivateChat(context.Background(), first, second)
			if err != nil {
				errors <- err
				return
			}
			results <- chat.ID
		}(i%2 == 0)
	}
	wg.Wait()
	close(results)
	close(errors)

	for err := range errors {
		t.Fatalf("CreatePrivateChat() error = %v", err)
	}
	var chatID uuid.UUID
	for result := range results {
		if chatID == uuid.Nil {
			chatID = result
		}
		if result != chatID {
			t.Fatalf("chat ID = %s, want %s", result, chatID)
		}
	}
	if len(chatRepo.chats) != 1 {
		t.Fatalf("chat count = %d, want 1", len(chatRepo.chats))
	}
	if len(membersRepo.members[chatID]) != 2 {
		t.Fatalf("member count = %d, want 2", len(membersRepo.members[chatID]))
	}
}
