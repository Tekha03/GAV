package memory

import (
	"context"
	"errors"
	"messenger/internal/model"
	"messenger/internal/repository"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrAttachmentNotFound = errors.New("attachment not found")
	ErrAttachmentExist = errors.New("attachment exist")
)

type AttachmentRepository struct {
	mu 			sync.RWMutex
	attachments map[uuid.UUID]*model.Attachment
}

func NewAttachmentRepository() repository.AttachmentRepository {
	return &AttachmentRepository{attachments: map[uuid.UUID]*model.Attachment{}}
}

func (ar *AttachmentRepository) Create(ctx context.Context, attachment *model.Attachment) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if attachment.ID != uuid.Nil {
		if _, found := ar.attachments[attachment.ID]; found {
			return ErrAttachmentExist
		}
	} else {
		attachment.ID = uuid.New()
	}

	ar.attachments[attachment.ID] = attachment
	return nil
}

func (ar *AttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Attachment, error) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	attachment, ok := ar.attachments[id]
	if !ok {
		return nil, ErrAttachmentNotFound
	}

	return attachment, nil
}

func (ar *AttachmentRepository) GetByMessage(ctx context.Context, messageID uuid.UUID) ([]*model.Attachment, error) {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	var result []*model.Attachment
	for _, attachment := range ar.attachments {
		if attachment.MessageID == messageID {
			result = append(result, attachment)
		}
	}

	return result, nil
}

func (ar *AttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if _, ok := ar.attachments[id]; !ok {
		return ErrAttachmentNotFound
	}

	delete(ar.attachments, id)
	return nil
}

func (ar *AttachmentRepository) DeleteByMessage(ctx context.Context, messageID uuid.UUID) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	for id, att := range ar.attachments {
		if att.MessageID == messageID {
			delete(ar.attachments, id)
		}
	}

	return nil
}
