package service

import (
	"messenger/internal/client"
	"messenger/internal/kafka"
	"messenger/internal/repository"
	"shared/events"
)

type ChatService struct {
	chatRepo       repository.ChatRepository
	membersRepo    repository.ChatMemberRepository
	messageRepo    repository.MessageRepository
	attachmentRepo repository.AttachmentRepository
	reactionRepo   repository.ReactionRepository
	pinnedRepo     repository.PinnedRepository
	typingRepo     repository.TypingRepository

	socialClient *client.SocialNetworkClient
	notClient    *client.NotificationClient

	producer kafka.EventProducer
}

func NewService(
	chatRepo repository.ChatRepository,
	membersRepo repository.ChatMemberRepository,
	messageRepo repository.MessageRepository,
	attachmentRepo repository.AttachmentRepository,
	reactionRepo repository.ReactionRepository,
	pinnedRepo repository.PinnedRepository,
	typingRepo repository.TypingRepository,

	socialClient *client.SocialNetworkClient,
	notClient *client.NotificationClient,

	producer kafka.EventProducer,

) Service {
	s := &ChatService{
		chatRepo:       chatRepo,
		membersRepo:    membersRepo,
		messageRepo:    messageRepo,
		attachmentRepo: attachmentRepo,
		reactionRepo:   reactionRepo,
		pinnedRepo:     pinnedRepo,
		typingRepo:     typingRepo,
		socialClient:   socialClient,
		notClient:      notClient,
		producer:       producer,
	}

	return s
}

func (s *ChatService) publishEvent(event events.Event) error {
	if s.producer == nil {
		return nil
	}
	return s.producer.PublishEvent(event)
}
