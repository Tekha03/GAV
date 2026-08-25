package service

import (
	"messenger/internal/client"
	"messenger/internal/kafka"
	"messenger/internal/repository"
	"shared/events"
)

type ChatService struct {
	transactionManager repository.TransactionManager
	chatRepo           repository.ChatRepository
	membersRepo        repository.ChatMemberRepository
	messageRepo        repository.MessageRepository
	attachmentRepo     repository.AttachmentRepository
	outboxRepo         repository.OutboxRepository
	reactionRepo       repository.ReactionRepository
	pinnedRepo         repository.PinnedRepository
	typingRepo         repository.TypingRepository

	socialClient *client.SocialNetworkClient
	notClient    *client.NotificationClient

	producer kafka.EventProducer
}

func NewService(
	transactionManager repository.TransactionManager,
	chatRepo repository.ChatRepository,
	membersRepo repository.ChatMemberRepository,
	messageRepo repository.MessageRepository,
	attachmentRepo repository.AttachmentRepository,
	outboxRepo repository.OutboxRepository,
	reactionRepo repository.ReactionRepository,
	pinnedRepo repository.PinnedRepository,
	typingRepo repository.TypingRepository,

	socialClient *client.SocialNetworkClient,
	notClient *client.NotificationClient,

	producer kafka.EventProducer,

) Service {
	s := &ChatService{
		transactionManager: transactionManager,
		chatRepo:           chatRepo,
		membersRepo:        membersRepo,
		messageRepo:        messageRepo,
		attachmentRepo:     attachmentRepo,
		outboxRepo:         outboxRepo,
		reactionRepo:       reactionRepo,
		pinnedRepo:         pinnedRepo,
		typingRepo:         typingRepo,
		socialClient:       socialClient,
		notClient:          notClient,
		producer:           producer,
	}

	return s
}

func (s *ChatService) publishEvent(event events.Event) error {
	if s.producer == nil {
		return nil
	}
	return s.producer.PublishEvent(event)
}
