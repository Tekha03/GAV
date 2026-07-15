package messengersync

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type ChatSnapshot struct {
	ID        uuid.UUID
	Members   map[uuid.UUID]struct{}
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MessageSnapshot struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Text      string
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Store struct {
	mu        sync.RWMutex
	chats     map[uuid.UUID]*ChatSnapshot
	messages  map[uuid.UUID]*MessageSnapshot
	reactions map[uuid.UUID]map[uuid.UUID]string
}

func NewStore() *Store {
	return &Store{
		chats:     make(map[uuid.UUID]*ChatSnapshot),
		messages:  make(map[uuid.UUID]*MessageSnapshot),
		reactions: make(map[uuid.UUID]map[uuid.UUID]string),
	}
}

func (s *Store) CreateChat(chatID uuid.UUID, members []uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat := &ChatSnapshot{
		ID:        chatID,
		Members:   make(map[uuid.UUID]struct{}, len(members)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	for _, memberID := range members {
		chat.Members[memberID] = struct{}{}
	}

	s.chats[chatID] = chat
}

func (s *Store) AddMember(chatID, userID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat, ok := s.chats[chatID]
	if !ok {
		chat = &ChatSnapshot{
			ID:      chatID,
			Members: make(map[uuid.UUID]struct{}),
		}
		s.chats[chatID] = chat
	}

	chat.Members[userID] = struct{}{}
	chat.UpdatedAt = time.Now()
}

func (s *Store) RemoveMember(chatID, userID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat, ok := s.chats[chatID]
	if !ok {
		return
	}

	delete(chat.Members, userID)
	chat.UpdatedAt = time.Now()
}

func (s *Store) DeleteChat(chatID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if chat, ok := s.chats[chatID]; ok {
		chat.Deleted = true
		chat.UpdatedAt = time.Now()
	}
}

func (s *Store) ChatMembers(chatID uuid.UUID) []uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chat, ok := s.chats[chatID]
	if !ok {
		return nil
	}

	members := make([]uuid.UUID, 0, len(chat.Members))
	for memberID := range chat.Members {
		members = append(members, memberID)
	}

	return members
}

func (s *Store) SaveMessage(messageID, chatID, senderID uuid.UUID, text string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages[messageID] = &MessageSnapshot{
		ID:        messageID,
		ChatID:    chatID,
		SenderID:  senderID,
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Store) EditMessage(messageID uuid.UUID, text string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	message, ok := s.messages[messageID]
	if !ok {
		return
	}

	message.Text = text
	message.UpdatedAt = time.Now()
}

func (s *Store) DeleteMessage(messageID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if message, ok := s.messages[messageID]; ok {
		message.Deleted = true
		message.UpdatedAt = time.Now()
	}
}

func (s *Store) Message(messageID uuid.UUID) (*MessageSnapshot, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	message, ok := s.messages[messageID]
	if !ok {
		return nil, false
	}

	copyMessage := *message
	return &copyMessage, ok
}

func (s *Store) AddReaction(messageID, userID uuid.UUID, emoji string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.reactions[messageID]; !ok {
		s.reactions[messageID] = make(map[uuid.UUID]string)
	}

	s.reactions[messageID][userID] = emoji
}

func (s *Store) RemoveReaction(messageID, userID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.reactions[messageID]; ok {
		delete(s.reactions[messageID], userID)
	}
}
