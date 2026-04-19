package kafka

import (
	"context"
	"shared/events"
)

type MessageUseCase interface {
	OnMessageSent(ctx context.Context, data events.MessageSentData)
	OnMessageEdited(ctx context.Context, data events.MessageEditedData)
	OnMessageDeleted(ctx context.Context, data events.MessageDeletedData)
}

type ChatUseCase interface {
	OnChatCreated(ctx context.Context, data events.ChatCreatedData)
	OnChatMemberAdded(ctx context.Context, data events.ChatMemberAddedData)
	OnChatMemberRemoved(ctx context.Context, data events.ChatMemberRemovedData)
	OnChatDeleted(ctx context.Context, data events.ChatDeletedData)
}

type ReactionUseCase interface {
	OnReactionAdded(ctx context.Context, data events.ReactionAddedData)
	OnReactionRemoved(ctx context.Context, data events.ReactionRemovedData)
}
