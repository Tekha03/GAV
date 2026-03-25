// messanger/chat/client/social_network_client.go
package client

import (
	"context"
	"google.golang.org/grpc"
	pb "api/social_network_gen"
	uuid "github.com/google/uuid"
)

type SocialNetworkClient struct {
	userClient    pb.UserServiceClient
	profileClient pb.ProfileServiceClient
	conn          *grpc.ClientConn
}

func NewSocialNetworkClient(addr string) (*SocialNetworkClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &SocialNetworkClient{
		userClient:    pb.NewUserServiceClient(conn),
		profileClient: pb.NewProfileServiceClient(conn),
		conn:          conn,
	}, nil
}

func (c *SocialNetworkClient) GetUser(ctx context.Context, id uuid.UUID) (*pb.UserResponse, error) {
	req := &pb.UserRequest{Id: id[:]}
	return c.userClient.GetUser(ctx, req)
}

func (c *SocialNetworkClient) GetUserProfile(ctx context.Context, userID uuid.UUID) (*pb.UserProfileResponse, error) {
	req := &pb.UserProfileRequest{UserId: userID[:]}
	return c.profileClient.GetUserProfile(ctx, req)
}
