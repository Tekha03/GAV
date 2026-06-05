package kafka

import (
	"encoding/json"
	"log"
	"shared/events"

	"github.com/IBM/sarama"
)

type Handler struct {
	message		MessageUseCase
	chat		ChatUseCase
	reaction	ReactionUseCase
}

func NewHandler(message MessageUseCase, chat ChatUseCase, reaction ReactionUseCase) (*Handler, error) {
	if message == nil {
		return nil, ErrMessageUseCaseNil
	}
	if chat == nil {
		return nil, ErrChatUseCaseNil
	}
	if reaction == nil {
		return nil, ErrReactionUseCaseNil
	}

	return &Handler{
		message:	message,
		chat:		chat,
		reaction:	reaction,
	}, nil
}


func (h *Handler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()

	for message := range claim.Messages() {
		var event events.Event
		if err := json.Unmarshal(message.Value, &event); err != nil {
			continue
		}

		switch event.EventType {
		case events.EventTypeMessageSent:
			var data events.MessageSentData
			_ = json.Unmarshal(event.Data, &data)
			h.message.OnMessageSent(ctx, data)
		case events.EventTypeMessageEdited:
			var data events.MessageEditedData
			_ = json.Unmarshal(event.Data, &data)
			h.message.OnMessageEdited(ctx, data)
		case events.EventTypeMessageDeleted:
			var data events.MessageDeletedData
			_ = json.Unmarshal(event.Data, &data)
			h.message.OnMessageDeleted(ctx, data)
		case events.EventTypeChatCreated:
			var data events.ChatCreatedData
			_ = json.Unmarshal(event.Data, &data)
			h.chat.OnChatCreated(ctx, data)
		case events.EventTypeChatMemberAdded:
			var data events.ChatMemberAddedData
			_ = json.Unmarshal(event.Data, &data)
			h.chat.OnChatMemberAdded(ctx, data)
		case events.EventTypeChatMemberRemoved:
			var data events.ChatMemberRemovedData
			_ = json.Unmarshal(event.Data, &data)
			h.chat.OnChatMemberRemoved(ctx, data)
		case events.EventTypeChatDeleted:
			var data events.ChatDeletedData
			_ = json.Unmarshal(event.Data, &data)
			h.chat.OnChatDeleted(ctx, data)
		case events.EventTypeReactionAdded:
			var data events.ReactionAddedData
			_ = json.Unmarshal(event.Data, &data)
			h.reaction.OnReactionAdded(ctx, data)
		case events.EventTypeReactionRemoved:
			var data events.ReactionRemovedData
			_ = json.Unmarshal(event.Data, &data)
			h.reaction.OnReactionRemoved(ctx, data)
		default:
			log.Println("unknown event type:", event.EventType)
		}

		session.MarkMessage(message, "")
	}

	return nil
}

func decode(input interface{}, output interface{}) {
	bytes, _ := json.Marshal(input)
	_ = json.Unmarshal(bytes, output)
}
