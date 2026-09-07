package container

import (
	"context"
	"database/sql"
	"errors"
	"messenger/internal/client"
	"messenger/internal/kafka"
	"messenger/internal/outbox"
	"messenger/internal/repository"
	"messenger/internal/service"
	orm "messenger/storage/gorm"
	rds "messenger/storage/redis"
	"messenger/transport/websocket"
	apperrors "shared/app_errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dependencyCheckTimeout = 5 * time.Second

type HybridContainer struct {
	websocket    *websocket.Hub
	gormRepo     *orm.Repository
	sqlDB        *sql.DB
	redis        *redis.Client
	socialClient *client.SocialNetworkClient
	notClient    *client.NotificationClient
	producer     *kafka.Producer
	outboxRepo   repository.OutboxRepository
}

func NewHybridContainer(
	ctx context.Context,
	postgresDSN string,
	redisAddr string,
	socialNetworkAddr string,
	producer *kafka.Producer,
) (*HybridContainer, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	websocketHub := websocket.NewHub()
	go websocketHub.Run()

	pgDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, apperrors.Wrap(apperrors.Internal, "failed to open postgres", err)
	}

	sqlDB, err := pgDB.DB()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.Internal, "failed to access postgres connection", err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	checkCtx, cancel := context.WithTimeout(ctx, dependencyCheckTimeout)
	defer cancel()

	if err := redisClient.Ping(checkCtx).Err(); err != nil {
		cause := errors.Join(err, redisClient.Close(), sqlDB.Close())
		return nil, apperrors.Wrap(apperrors.ServiceUnavailable, "failed to connect to redis", cause)
	}

	socialClient, err := client.NewSocialNetworkClient(socialNetworkAddr)
	if err != nil {
		cause := errors.Join(err, redisClient.Close(), sqlDB.Close())
		return nil, apperrors.Wrap(apperrors.Internal, "failed to create social network client", cause)
	}

	notClient, err := client.NewNotificationClient(socialNetworkAddr)
	if err != nil {
		cause := errors.Join(err, socialClient.Close(), redisClient.Close(), sqlDB.Close())
		return nil, apperrors.Wrap(apperrors.Internal, "failed to create notification client", cause)
	}

	gormRepo := orm.NewRepository(pgDB)

	return &HybridContainer{
		websocket:    websocketHub,
		gormRepo:     gormRepo,
		sqlDB:        sqlDB,
		redis:        redisClient,
		socialClient: socialClient,
		notClient:    notClient,
		producer:     producer,
		outboxRepo:   orm.NewOutboxRepository(gormRepo),
	}, nil
}

func (c *HybridContainer) ChatService() service.Service {

	return service.NewService(
		c.gormRepo,
		orm.NewChatRepository(c.gormRepo),
		orm.NewChatMemberRepository(c.gormRepo),
		orm.NewMessageRepository(c.gormRepo),
		orm.NewAttachmentRepository(c.gormRepo),
		c.outboxRepo,
		orm.NewReactionRepository(c.gormRepo),
		rds.NewPinnedRepository(c.redis),
		rds.NewTypingRepository(c.redis),
		c.socialClient,
		c.notClient,
		c.producer,
		c.websocket,
	)
}

func (c *HybridContainer) RunOutboxWorker(ctx context.Context) {
	if c == nil {
		return
	}
	outbox.NewWorker(c.outboxRepo, c.producer).Run(ctx)
}

func (c *HybridContainer) Close() error {
	if c == nil {
		return nil
	}

	var closeErrors []error
	closeErrors = appendIfError(closeErrors, closeNotificationClient(c.notClient))
	closeErrors = appendIfError(closeErrors, closeSocialClient(c.socialClient))
	closeErrors = appendIfError(closeErrors, closeRedis(c.redis))
	closeErrors = appendIfError(closeErrors, closeSQL(c.sqlDB))

	if err := errors.Join(closeErrors...); err != nil {
		return apperrors.Wrap(apperrors.Internal, "failed to close container resources", err)
	}

	return nil
}

func appendIfError(errs []error, err error) []error {
	if err != nil {
		return append(errs, err)
	}
	return errs
}

func closeNotificationClient(client *client.NotificationClient) error {
	if client == nil {
		return nil
	}
	return client.Close()
}

func closeSocialClient(client *client.SocialNetworkClient) error {
	if client == nil {
		return nil
	}
	return client.Close()
}

func closeRedis(client *redis.Client) error {
	if client == nil {
		return nil
	}
	return client.Close()
}

func closeSQL(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}

func (c *HybridContainer) WebSocketHub() *websocket.Hub {
	return c.websocket
}
