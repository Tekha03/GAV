// messanger/chat/container/gorm_container.go
package container

import (
	"messanger/chat/service"
	orm "messanger/storage/gorm"
	rds "messanger/storage/redis"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
type HybridContainer struct {
    gormRepo *orm.Repository
    redis    *redis.Client
}

func NewHybridContainer(postgresDSN, redisAddr string) (*HybridContainer, error) {
    pgDB, _ := gorm.Open(postgres.Open(postgresDSN))
    gormRepo := orm.NewRepository(pgDB)
    
    redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
    
    return &HybridContainer{
        gormRepo: gormRepo,
        redis:    redisClient,
    }, nil
}

func (c *HybridContainer) ChatService() service.Service {
    return service.NewService(
        orm.NewChatRepository(c.gormRepo),
        orm.NewChatMemberRepository(c.gormRepo),
        orm.NewMessageRepository(c.gormRepo),
        orm.NewAttachmentRepository(c.gormRepo),
        orm.NewReactionRepository(c.gormRepo),
        rds.NewPinnedRepository(c.redis),
        rds.NewTypingRepository(c.redis),
    )
}
