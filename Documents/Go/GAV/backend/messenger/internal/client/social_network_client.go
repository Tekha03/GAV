// messanger/chat/client/social_network_client.go
package client

import (
	"context"
	"google.golang.org/grpc"
	pb "api/gen/social/v1"
	uuid "github.com/google/uuid"
)

type SocialNetworkClient struct {
	socialClient pb.SocialServiceClient
	conn          *grpc.ClientConn
}

func NewSocialNetworkClient(addr string) (*SocialNetworkClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &SocialNetworkClient{
		socialClient: pb.NewSocialServiceClient(conn),
		conn:          conn,
	}, nil
}

func (c *SocialNetworkClient) GetUserProfile(ctx context.Context, userID uuid.UUID) (*pb.UserProfile, error) {
	resp, err := c.socialClient.GetProfile(ctx, &pb.GetProfileRequest{UserId: userID.String()})
	return resp, err
  }

  func (c *SocialNetworkClient) Login(ctx context.Context, email, password string) (*pb.LoginResponse, error) {
	// return c.socialClient.Login(ctx, &pb.LoginRequest{Email: email, Password: password})
	return nil, nil
  }
