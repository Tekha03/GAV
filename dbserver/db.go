package dbserver

import (
	"fmt"
	"log/slog"

	"gav/internal/comment"
	"gav/internal/dog"
	"gav/internal/follow"
	"gav/internal/like"
	"gav/internal/post"
	"gav/internal/profile"
	"gav/internal/settings"
	"gav/internal/stats"
	"gav/internal/token"
	"gav/internal/user"
	"gav/internal/vaccination"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(path string, logger *slog.Logger) (*gorm.DB, error) {
	if path == "" {
		path = "social.db"
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite: %w", err)
	}

	logger.Info("database opened", "path", path)
	sqlDB, _ := db.DB()
	sqlDB.Exec("PRAGMA foreign_keys = ON;")

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
		// posts
		"CREATE INDEX IF NOT EXISTS idx_posts_user_id_created ON posts(user_id, created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC)",

		// comments
		"CREATE INDEX IF NOT EXISTS idx_comments_post_id_created ON comments(post_id, created_at ASC)",

		// likes
		"CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id)",
		"CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id)",

		// follows
		"CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows(follower_id)",
		"CREATE INDEX IF NOT EXISTS idx_follows_following_id ON follows(following_id)",

		// dogs
		"CREATE INDEX IF NOT EXISTS idx_dogs_owner_id ON dogs(owner_id)",

		// vaccinations
		"CREATE INDEX IF NOT EXISTS idx_vaccinations_dog_id ON vaccinations(dog_id)",

		// refresh-tokens
		"CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash)",
		"CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at)",
	}

	for _, idx := range indexes {
		if err := db.Exec(idx); err != nil {
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
