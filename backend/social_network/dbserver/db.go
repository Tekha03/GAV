package dbserver

import (
	"fmt"
	"log/slog"
	"strings"

	"social_network/internal/comment"
	"social_network/internal/dog"
	"social_network/internal/follow"
	"social_network/internal/like"
	"social_network/internal/post"
	"social_network/internal/profile"
	"social_network/internal/settings"
	"social_network/internal/stats"
	"social_network/internal/token"
	"social_network/internal/user"
	"social_network/internal/vaccination"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(driver, path, postgresDSN string, logger *slog.Logger) (*gorm.DB, error) {
	driver = strings.ToLower(strings.TrimSpace(driver))
	if driver == "" {
		driver = "sqlite"
	}

	var dialector gorm.Dialector
	var source string

	switch driver {
	case "sqlite":
		if path == "" {
			path = "social.db"
		}
		dialector = sqlite.Open(path)
		source = path
	case "postgres", "postgresql":
		if postgresDSN == "" {
			return nil, fmt.Errorf("postgres dsn is empty")
		}
		dialector = postgres.Open(postgresDSN)
		source = "postgres"
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedDriver, driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", driver, err)
	}

	logger.Info("database opened", "driver", driver, "source", source)
	if driver == "sqlite" {
		sqlDB, _ := db.DB()
		sqlDB.Exec("PRAGMA foreign_keys = ON;")
	}

	models := []interface{}{
		&user.User{},
		&profile.UserProfile{},
		&settings.UserSettings{},
		&post.Post{},
		&comment.Comment{},
		&like.Like{},
		&follow.Follow{},
		&dog.Dog{},
		&vaccination.Vaccination{},
		&token.RefreshToken{},
		&stats.UserStats{},
		&stats.PostStats{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

	logger.Info("auto migration completed", "models_count", len(models))

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_posts_user_id_created ON posts(user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_comments_post_id_created ON comments(post_id, created_at ASC)",
		"CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id)",
		"CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows(follower_id)",
		"CREATE INDEX IF NOT EXISTS idx_follows_following_id ON follows(following_id)",
		"CREATE INDEX IF NOT EXISTS idx_dogs_owner_id ON dogs(owner_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_location_status_visibility ON users(location_status, visibility)",
		"CREATE INDEX IF NOT EXISTS idx_vaccinations_dog_id ON vaccinations(dog_id)",
		"CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash)",
		"CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at)",
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			logger.Warn("failed to create index", "sql", idx, "error", err)
		}
	}

	logger.Info("additional indexes created")

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
