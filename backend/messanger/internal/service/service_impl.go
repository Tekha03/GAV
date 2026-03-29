package service

import (
	"messanger/internal/client"
	"messanger/internal/repository"
)

type ChatService struct {
	chatRepo      repository.ChatRepository
	membersRepo   repository.ChatMemberRepository
	messageRepo   repository.MessageRepository
	attachmentRepo repository.AttachmentRepository
	reactionRepo  repository.ReactionRepository
    pinnedRepo    repository.PinnedRepository
    typingRepo repository.TypingRepository

    socialClient *client.SocialNetworkClient
	notClient    *client.NotificationClient
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