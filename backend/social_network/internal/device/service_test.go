package device

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type repositoryStub struct {
	savedUser    uuid.UUID
	savedToken   string
	deletedUser  uuid.UUID
	deletedToken string
	err          error
}

func (r *repositoryStub) Create(context.Context, *DeviceToken) error { return r.err }
func (r *repositoryStub) GetByUser(context.Context, uuid.UUID) ([]*DeviceToken, error) {
	return nil, r.err
}
func (r *repositoryStub) SaveForUser(_ context.Context, userID uuid.UUID, token string) error {
	r.savedUser, r.savedToken = userID, token
	return r.err
}
func (r *repositoryStub) DeleteForUser(_ context.Context, userID uuid.UUID, token string) error {
	r.deletedUser, r.deletedToken = userID, token
	return r.err
}

func TestServiceRegisterTrimsAndPersistsToken(t *testing.T) {
	repo := &repositoryStub{}
	service := NewService(repo)
	userID := uuid.New()

	if err := service.Register(context.Background(), userID, "  device-token  "); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if repo.savedUser != userID || repo.savedToken != "device-token" {
		t.Fatalf("saved token = (%s, %q)", repo.savedUser, repo.savedToken)
	}
}

func TestServiceRejectsInvalidTokens(t *testing.T) {
	service := NewService(&repositoryStub{})
	for _, test := range []struct {
		name   string
		userID uuid.UUID
		token  string
	}{
		{name: "missing user", userID: uuid.Nil, token: "token"},
		{name: "empty token", userID: uuid.New(), token: "  "},
		{name: "oversized token", userID: uuid.New(), token: strings.Repeat("x", 4097)},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := service.Register(context.Background(), test.userID, test.token); !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("Register() error = %v, want ErrInvalidToken", err)
			}
			if err := service.Unregister(context.Background(), test.userID, test.token); !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("Unregister() error = %v, want ErrInvalidToken", err)
			}
		})
	}
}

func TestServiceUnregisterPersistsAndPropagatesErrors(t *testing.T) {
	wantErr := errors.New("storage unavailable")
	repo := &repositoryStub{err: wantErr}
	service := NewService(repo)
	userID := uuid.New()

	err := service.Unregister(context.Background(), userID, " token ")
	if !errors.Is(err, wantErr) {
		t.Fatalf("Unregister() error = %v, want %v", err, wantErr)
	}
	if repo.deletedUser != userID || repo.deletedToken != "token" {
		t.Fatalf("deleted token = (%s, %q)", repo.deletedUser, repo.deletedToken)
	}
}
