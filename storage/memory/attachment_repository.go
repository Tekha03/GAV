package memory

import (
	"context"
	"errors"
	"gav/internal/chat"
	"sync"
)

var (
	ErrAttachmentNotFound = errors.New("attachment not found")
)

type AttachmentRepository struct {
	mu 			sync.RWMutex
	lastID		uint
	attachments map[uint]*chat.Attachment
}

func NewAttachmentrepository() *AttachmentRepository {
	return &AttachmentRepository{attachments: map[uint]*chat.Attachment{}}
}

func (ar *AttachmentRepository) Create(ctx context.Context, attachment *chat.Attachment) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.lastID++
	attachment.ID = ar.lastID

	ar.attachments[attachment.ID] = attachment
	return nil
}

func (ar *AttachmentRepository) GetByID(ctx context.Context, id uint) (*chat.Attachment, error) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	attachment, ok := ar.attachments[id]
	if !ok {
		return nil, ErrAttachmentNotFound
	}

	return attachment, nil
}

func (ar *AttachmentRepository) GetByMessage(ctx context.Context, messageID uint) ([]*chat.Attachment, error) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	var result []*chat.Attachment
	for _, attachment := range ar.attachments {
		if attachment.MessageID == messageID {
			result = append(result, attachment)
		}
	}

	return result, nil
}

func (ar *AttachmentRepository) Delete(ctx context.Context, id uint) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if _, ok := ar.attachments[id]; !ok {
		return ErrAttachmentNotFound
	}

	delete(ar.attachments, id)
	return nil
}

func (ar *AttachmentRepository) DeleteByMessage(ctx context.Context, messageID uint) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	for id, att := range ar.attachments {
		if att.MessageID == messageID {
			delete(ar.attachments, id)
		}
	}

	return nil
}