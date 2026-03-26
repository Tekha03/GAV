// messanger/chat/container/gorm_container.go
package container

import (
	"messanger/chat/client"
	"messanger/chat/service"
	orm "messanger/storage/gorm"
	rds "messanger/storage/redis"
	"messanger/transport/websocket"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
type HybridContainer struct {
    gormRepo *orm.Repository
    redis    *redis.Client
    grpcConn *grpc.ClientConn 
    socialClient *client.SocialNetworkClient
	notClient    *client.NotificationClient
}

func NewHybridContainer(postgresDSN, redisAddr, socialNetworkAddr string) (*HybridContainer, error) {
	pgDB, err := gorm.Open(postgres.Open(postgresDSN))
	if err != nil {
		return nil, err
	}
	gormRepo := orm.NewRepository(pgDB)

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})

	grpcConn, err := grpc.Dial(socialNetworkAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	socialClient, err := client.NewSocialNetworkClient(socialNetworkAddr)
	if err != nil {
		grpcConn.Close()
		return nil, err
	}

	notClient, err := client.NewNotificationClient(socialNetworkAddr)
	if err != nil {
		grpcConn.Close()
		return nil, err
	}

	return &HybridContainer{
		gormRepo:     gormRepo,
		redis:        redisClient,
		grpcConn:     grpcConn,
		socialClient: socialClient,
		notClient:    notClient,
	}, nil
}

func (c *HybridContainer) ChatService() service.Service {
    websocketHub := websocket.NewHub()
	go websocketHub.Run()

	return service.NewService(
		orm.NewChatRepository(c.gormRepo),
		orm.NewChatMemberRepository(c.gormRepo),
		orm.NewMessageRepository(c.gormRepo),
		orm.NewAttachmentRepository(c.gormRepo),
		orm.NewReactionRepository(c.gormRepo),
		rds.NewPinnedRepository(c.redis),
		rds.NewTypingRepository(c.redis),

		c.socialClient,
		c.notClient,
	)
}

func (c *HybridContainer) Close() error {
	if c.grpcConn != nil {
		return c.grpcConn.Close()
	}
	return nil
}
