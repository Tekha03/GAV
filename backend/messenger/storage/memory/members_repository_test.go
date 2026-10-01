package memory

import (
	"context"
	"testing"

	"messenger/internal/model"

	"github.com/google/uuid"
)

func TestMembersRepositoryLifecycleAndQueries(t *testing.T) {
	ctx := context.Background()
	repo := NeMembersRepository()
	chatID, secondChatID := uuid.New(), uuid.New()
	firstUser, secondUser, thirdUser := uuid.New(), uuid.New(), uuid.New()

	first := &model.ChatMember{ChatID: chatID, UserID: firstUser, Role: model.Member}
	second := &model.ChatMember{ChatID: chatID, UserID: secondUser, Role: model.Member}
	for _, member := range []*model.ChatMember{
		first,
		second,
		{ChatID: secondChatID, UserID: firstUser, Role: model.Member},
		{ChatID: secondChatID, UserID: thirdUser, Role: model.Member},
	} {
		if err := repo.AddMember(ctx, member); err != nil {
			t.Fatal(err)
		}
	}
	if err := repo.AddMember(ctx, first); err == nil {
		t.Fatal("duplicate AddMember() error = nil")
	}

	members, err := repo.GetMembers(ctx, chatID)
	if err != nil || len(members) != 2 {
		t.Fatalf("GetMembers() count=%d error=%v", len(members), err)
	}
	if missing, err := repo.GetMembers(ctx, uuid.New()); err != nil || len(missing) != 0 {
		t.Fatalf("GetMembers() missing = %+v, %v", missing, err)
	}
	if exists, err := repo.MemberExists(ctx, firstUser, chatID); err != nil || !exists {
		t.Fatalf("MemberExists() = %v, %v", exists, err)
	}
	if exists, err := repo.MemberExists(ctx, thirdUser, chatID); err != nil || exists {
		t.Fatalf("MemberExists(non-member) = %v, %v", exists, err)
	}

	admin := model.Admin
	if err := repo.UpdateRole(ctx, chatID, firstUser, &admin); err != nil {
		t.Fatal(err)
	}
	role, err := repo.GetRole(ctx, firstUser, chatID)
	if err != nil || role == nil || *role != model.Admin {
		t.Fatalf("GetRole() = %v, %v", role, err)
	}
	if err := repo.SetMuted(ctx, chatID, firstUser, true); err != nil || !first.Muted {
		t.Fatalf("SetMuted() member=%+v error=%v", first, err)
	}

	messageID := uuid.New()
	if err := repo.UpdateLastReadMessageID(ctx, chatID, firstUser, messageID); err != nil {
		t.Fatal(err)
	}
	lastRead, err := repo.GetLastReadMessageID(ctx, chatID, firstUser)
	if err != nil || lastRead != messageID {
		t.Fatalf("GetLastReadMessageID() = %s, %v", lastRead, err)
	}

	privateChatID, err := repo.FindPrivateChatBetween(ctx, firstUser, secondUser)
	if err != nil || privateChatID != chatID {
		t.Fatalf("FindPrivateChatBetween() = %s, %v", privateChatID, err)
	}
	if missingID, err := repo.FindPrivateChatBetween(ctx, secondUser, thirdUser); err != nil || missingID != uuid.Nil {
		t.Fatalf("FindPrivateChatBetween() missing = %s, %v", missingID, err)
	}
	chats, err := repo.GetUserChats(ctx, firstUser)
	if err != nil || len(chats) != 2 {
		t.Fatalf("GetUserChats() = %+v, %v", chats, err)
	}

	if err := repo.RemoveMember(ctx, secondUser, chatID); err != nil {
		t.Fatal(err)
	}
	if err := repo.RemoveMember(ctx, secondUser, chatID); err == nil {
		t.Fatal("RemoveMember() missing member error = nil")
	}
}

func TestMembersRepositoryMissingMemberMutationsFail(t *testing.T) {
	ctx := context.Background()
	repo := NeMembersRepository()
	chatID, userID := uuid.New(), uuid.New()
	role := model.Admin

	if err := repo.RemoveMember(ctx, userID, chatID); err == nil {
		t.Fatal("RemoveMember() missing chat error = nil")
	}
	if err := repo.UpdateRole(ctx, chatID, userID, &role); err == nil {
		t.Fatal("UpdateRole() missing chat error = nil")
	}
	if err := repo.SetMuted(ctx, chatID, userID, true); err == nil {
		t.Fatal("SetMuted() missing chat error = nil")
	}
	if _, err := repo.GetRole(ctx, userID, chatID); err == nil {
		t.Fatal("GetRole() missing chat error = nil")
	}
	if _, err := repo.GetLastReadMessageID(ctx, chatID, userID); err == nil {
		t.Fatal("GetLastReadMessageID() missing chat error = nil")
	}
	if err := repo.UpdateLastReadMessageID(ctx, chatID, userID, uuid.New()); err == nil {
		t.Fatal("UpdateLastReadMessageID() missing chat error = nil")
	}
}
