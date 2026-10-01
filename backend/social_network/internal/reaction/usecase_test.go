package reaction

import (
	"context"
	"testing"

	"shared/events"
	messengersync "social_network/internal/messenger_sync"

	"github.com/google/uuid"
)

func TestUseCaseReactionLifecycleWithoutNotifier(t *testing.T) {
	store := messengersync.NewStore()
	useCase := NewUseCase(store, nil)
	messageID, userID := uuid.New(), uuid.New()

	useCase.OnReactionAdded(context.Background(), events.ReactionAddedData{
		MessageID: messageID, UserID: userID, Reaction: "🔥",
	})
	useCase.OnReactionRemoved(context.Background(), events.ReactionRemovedData{
		MessageID: messageID, UserID: userID,
	})
}
