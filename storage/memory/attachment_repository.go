package memory

import (
	"gav/internal/chat"
	"sync"
)

type AttachmentRepository struct {
	mu 			sync.RWMutex
	lastID		uint
	attachments map[int]*chat.Attachment
}

func NewAttachmentrepository() *AttachmentRepository {
	return &AttachmentRepository{attachments: map[int]*chat.Attachment{}}
}