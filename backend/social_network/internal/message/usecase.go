package message

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

func (u *UseCase) OnMessageSent(ctx context.Context, data events.MessageSentData) {
	log.Println("handle message sent:", data.Text)
}

func (u *UseCase) OnMessageEdited(ctx context.Context, data events.MessageEditedData) {
	log.Println("handle message edited:", data.Text)
}

func (u *UseCase) OnMessageDeleted(ctx context.Context, data events.MessageDeletedData) {
	log.Println("handle message deleted:", data.MessageID)
}
