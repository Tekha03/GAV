package postgres

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) (*BaseRepository, error) {
	if db == nil {
		return nil, ErrDBNil
	}

	return &BaseRepository{db: db}, nil
}

func (r *BaseRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}
