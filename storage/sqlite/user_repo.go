package sqlite

import (
	"context"
	"errors"

	"gav/internal/user"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists = errors.New("user already exists")
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(ctx context.Context, u *user.User) error {
	var existing user.User
	err := ur.db.WithContext(ctx).Where("settings_email = ?", u.Email).First(&existing).Error

	if err == nil {
		return ErrUserExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return ur.db.WithContext(ctx).Create(u).Error
}

func (ur *UserRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	if err := ur.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := ur.db.WithContext(ctx).Where("settings_email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (ur *UserRepository) Update(ctx context.Context, u *user.User) error {
	updated := ur.db.WithContext(ctx).Save(u)
	if updated.Error != nil {
		return updated.Error
	}

	if updated.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (ur *UserRepository) Delete(ctx context.Context, id uint) error {
	deleted := ur.db.WithContext(ctx).Delete(&user.User{}, "id = ?", id)
	if deleted.Error != nil {
		return deleted.Error
	}

	if deleted.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
