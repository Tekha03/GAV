package reaction

import (
	"context"
	"log"
	"shared/events"
)

type UseCase struct{
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) OnReactionAdded(ctx context.Context, data events.ReactionAddedData) {
	log.Println("handle reaction added:", data.MessageID, data.Reaction)
}

func (u *UseCase) OnReactionRemoved(ctx context.Context, data events.ReactionRemovedData) {
	log.Println("handle reaction removed:", data.MessageID, data.Reaction)
}
