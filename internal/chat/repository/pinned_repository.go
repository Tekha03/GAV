package repository

import "github.com/google/uuid"

type PinnedRepository interface {
	Pin(chatID, messageID uuid.UUID) error
	Unpin(chatID, messageID uuid.UUID) error
	GetByChatID(chatID uuid.UUID) []uuid.UUID
}