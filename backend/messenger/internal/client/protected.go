package client

import (
	authv1 "api/gen/auth/v1"
	socialv1 "api/gen/social/v1"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sony/gobreaker/v2"
)

type ProtectedNotificationClient struct {
	inner   NotifiClient
	breaker *gobreaker.CircuitBreaker[struct{}]
}

func NewProtectedNotificationClient(inner NotifiClient) NotifiClient {
	breaker := gobreaker.NewCircuitBreaker[struct{}](gobreaker.Settings{
		Name:        "notification-service",
		MaxRequests: 2,
		Interval:    30 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	})

	return &ProtectedNotificationClient{
		inner:   inner,
		breaker: breaker,
	}
}

func (c *ProtectedNotificationClient) SendNewMessage(
	ctx context.Context,
	receiverID, senderID uuid.UUID,
	senderName string,
	text string,
	chatID string,
) error {
	_, err := c.breaker.Execute(func() (struct{}, error) {
		return struct{}{}, c.inner.SendNewMessage(ctx, receiverID, senderID, senderName, text, chatID)
	})

	return err
}

func (c *ProtectedNotificationClient) Close() error {
	if c == nil || c.inner == nil {
		return nil
	}
	return c.inner.Close()
}

type ProtectedSocialClient struct {
	inner SocialClient

	profileBreaker *gobreaker.CircuitBreaker[*socialv1.UserProfile]
	loginBreaker   *gobreaker.CircuitBreaker[*authv1.LoginResponse]
}

func NewProtectedSocialClient(inner SocialClient) SocialClient {
	profileBreaker := gobreaker.NewCircuitBreaker[*socialv1.UserProfile](gobreaker.Settings{
		Name:        "social-network.get-user-profile",
		MaxRequests: 2,
		Interval:    30 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	})

	loginBreaker := gobreaker.NewCircuitBreaker[*authv1.LoginResponse](gobreaker.Settings{
		Name:        "social-network.login",
		MaxRequests: 2,
		Interval:    30 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	})

	return &ProtectedSocialClient{
		inner:          inner,
		profileBreaker: profileBreaker,
		loginBreaker:   loginBreaker,
	}
}

func (c *ProtectedSocialClient) GetUserProfile(ctx context.Context, userID uuid.UUID) (*socialv1.UserProfile, error) {
	return c.profileBreaker.Execute(func() (*socialv1.UserProfile, error) {
		return c.inner.GetUserProfile(ctx, userID)
	})
}

func (c *ProtectedSocialClient) Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error) {
	return c.loginBreaker.Execute(func() (*authv1.LoginResponse, error) {
		return c.inner.Login(ctx, email, password)
	})
}

func (c *ProtectedSocialClient) Close() error {
	if c == nil || c.inner == nil {
		return nil
	}
	return c.inner.Close()
}
