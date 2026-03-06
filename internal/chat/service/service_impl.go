package service

import (
	"gav/internal/chat/repository"
)

type ChatService struct {
	chatRepo      repository.ChatRepository
	membersRepo   repository.ChatMemberRepository
	messageRepo   repository.MessageRepository
	attachmentRepo repository.AttachmentRepository
	reactionRepo  repository.ReactionRepository
    pinnedRepo    repository.PinnedRepository
    typingRepo repository.TypingRepository
}

func NewService(
    chatRepo repository.ChatRepository,
    membersRepo repository.ChatMemberRepository,
    messageRepo repository.MessageRepository,
    attachmentRepo repository.AttachmentRepository,
    reactionRepo repository.ReactionRepository,
    pinnedRepo repository.PinnedRepository,
    typingRepo repository.TypingRepository,
) Service {
    return &ChatService{
        chatRepo:      chatRepo,
        membersRepo:   membersRepo,
        messageRepo:   messageRepo,
        attachmentRepo: attachmentRepo,
        reactionRepo:  reactionRepo,
        pinnedRepo:    pinnedRepo,
        typingRepo:    typingRepo,
    }
}