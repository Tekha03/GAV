package firebase

import (
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Client struct {
	messaging *messaging.Client
}

func NewClient(serviceAccountKeyPath string) (*Client, error) {
	ctx := context.Background()
	conf := &firebase.Config{}

	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(ctx)
	return &Client{
		messaging: client,
	}, err
}

func (c *Client) SendPush(
	ctx context.Context,
	token string,
	title, body string,
	data map[string]string,
) error {
	if c == nil || c.messaging == nil {
		return nil
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
