package service

import (
    "messenger/internal/kafka"
	"messenger/internal/client"
	"messenger/internal/repository"
)

type ChatService struct {
	chatRepo        repository.ChatRepository
	membersRepo     repository.ChatMemberRepository
	messageRepo     repository.MessageRepository
	attachmentRepo  repository.AttachmentRepository
	reactionRepo    repository.ReactionRepository
    pinnedRepo      repository.PinnedRepository
    typingRepo      repository.TypingRepository

    socialClient    *client.SocialNetworkClient
	notClient       *client.NotificationClient

    producer        kafka.Producer
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
	notClient    *client.NotificationClient,

) Service {
    s := &ChatService{
        chatRepo:      chatRepo,
        membersRepo:   membersRepo,
        messageRepo:   messageRepo,
        attachmentRepo: attachmentRepo,
        reactionRepo:  reactionRepo,
        pinnedRepo:    pinnedRepo,
        typingRepo:    typingRepo,
        socialClient: socialClient,
        notClient: notClient,
    }

    return s
}
