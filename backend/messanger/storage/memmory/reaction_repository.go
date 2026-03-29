package memory

import (
	"context"
	"messanger/internal/model"
	"messanger/internal/repository"
	"sync"

	"github.com/google/uuid"
)

type ReactionRepository struct {
	mu 			sync.RWMutex
	reactions 	map[uuid.UUID]map[uuid.UUID]*model.Reaction
}

func NewReactionRepository() repository.ReactionRepository {
	return &ReactionRepository{reactions: make(map[uuid.UUID]map[uuid.UUID]*model.Reaction)}
}

func (rr *ReactionRepository) Add(ctx context.Context, reaction *model.Reaction) error {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if _, ok := rr.reactions[reaction.MessageID]; !ok {
		rr.reactions[reaction.MessageID] = make(map[uuid.UUID]*model.Reaction)
	}

	if _, exists := rr.reactions[reaction.MessageID][reaction.UserID]; exists {
		return repository.ErrReactionExists
	}

	if reaction.ID == uuid.Nil {
		reaction.ID = uuid.New()
	}

	rr.reactions[reaction.MessageID][reaction.UserID] = reaction
	return nil
}

func (rr *ReactionRepository) Remove(ctx context.Context, messageID, userID uuid.UUID) error {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if _, ok := rr.reactions[messageID]; !ok {
		return repository.ErrReactionNotFound
	}

	if _, exists := rr.reactions[messageID][userID]; !exists {
		return repository.ErrReactionNotFound
	}

	delete(rr.reactions[messageID], userID)
	return nil
}