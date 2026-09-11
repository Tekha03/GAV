// messanger/chat/client/social_network_client.go
package client

import (
	authv1 "api/gen/auth/v1"
	socialv1 "api/gen/social/v1"
	"context"
	apperrors "shared/app_errors"
	"shared/retry"

	uuid "github.com/google/uuid"
	"google.golang.org/grpc"
)

type SocialNetworkClient struct {
	socialClient socialv1.SocialServiceClient
	authClient   authv1.AuthServiceClient
	conn         *grpc.ClientConn
}

func (c *SocialNetworkClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return apperrors.Wrap(apperrors.Internal, "failed to close social network client", c.conn.Close())
}

func NewSocialNetworkClient(addr string) (*SocialNetworkClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ServiceUnavailable, "failed to create social network client", err)
	}
	return &SocialNetworkClient{
		socialClient: socialv1.NewSocialServiceClient(conn),
		authClient:   authv1.NewAuthServiceClient(conn),
		conn:         conn,
	}, nil
}

func (c *SocialNetworkClient) GetUserProfile(ctx context.Context, userID uuid.UUID) (*socialv1.UserProfile, error) {
	resp, err := retry.DoValue(ctx, retry.DefaultConfig(), func(ctx context.Context) (*socialv1.UserProfile, error) {
		return c.getUserProfileOnce(ctx, userID)
	})
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ServiceUnavailable, "failed to get user profile", err)
	}
	return resp, nil
}

func (c *SocialNetworkClient) getUserProfileOnce(ctx context.Context, userID uuid.UUID) (*socialv1.UserProfile, error) {
	resp, err := c.socialClient.GetProfile(ctx, &socialv1.GetProfileRequest{UserId: userID.String()})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *SocialNetworkClient) Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error) {
	resp, err := c.authClient.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, apperrors.Wrap(apperrors.ServiceUnavailable, "failed to log in through auth service", err)
	}
	return resp, nil
}
