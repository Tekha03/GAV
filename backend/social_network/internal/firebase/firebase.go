package firebase

import (
	"context"
	"errors"
	"fmt"
	"os"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Client struct {
	messaging *messaging.Client
}

var ErrDisabled = errors.New("Firebase Cloud Messaging is disabled")

func NewClient(ctx context.Context, serviceAccountKeyPath string) (*Client, error) {
	if serviceAccountKeyPath == "" {
		return nil, errors.New("firebase credentials file is required")
	}
	if _, err := os.Stat(serviceAccountKeyPath); err != nil {
		return nil, fmt.Errorf("firebase credentials file unavailable: %w", err)
	}

	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, &firebase.Config{}, opt)
	if err != nil {
		return nil, fmt.Errorf("initialize firebase: %w", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize firebase messaging: %w", err)
	}
	return &Client{messaging: client}, nil
}

func (c *Client) SendPush(
	ctx context.Context,
	token string,
	title, body string,
	data map[string]string,
) error {
	if c == nil || c.messaging == nil {
		return ErrDisabled
	}

	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	_, err := c.messaging.Send(ctx, msg)
	return err
}
