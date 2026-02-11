package chat

import "context"

type MessageRepository interface {
	Create(ctx context.Context, message *Message) error 
	GetMessages(ctx context.Context, chatID, BeforeID uint, limit int) ([]*Message, error)
	
}