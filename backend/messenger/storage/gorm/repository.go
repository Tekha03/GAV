package gorm

import (
	"context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type transactionContextKey struct{}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionContextKey{}).(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) WithPrivateChatLock(ctx context.Context, privateKey string, fn func(context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtextextended(?, 0))", privateKey).Error; err != nil {
			return internalError("failed to lock private chat", err)
		}
		return fn(context.WithValue(ctx, transactionContextKey{}, tx))
	})
}
