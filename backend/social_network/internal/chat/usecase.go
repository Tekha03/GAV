package chat

import (
	"context"
	"log"
	"shared/events"
)

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) OnChatCreated(ctx context.Context, data events.ChatCreatedData) {
	log.Println("handle chat created:", data.ChatID)
}

func (u *UseCase) OnChatMemberAdded(ctx context.Context, data events.ChatMemberAddedData) {
	log.Println("handle chat member added:", data.ChatID, data.UserID)
}

func (u *UseCase) OnChatMemberRemoved(ctx context.Context, data events.ChatMemberRemovedData) {
	log.Println("handle chat member removed:", data.ChatID, data.UserID)
}

func (u *UseCase) OnChatDeleted(ctx context.Context, data events.ChatDeletedData) {
	log.Println("handle chat deleted:", data.ChatID)
}
