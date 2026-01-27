package sqlite

import (
	"gav/internal/role"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (rr *RoleRepository) GetByName(name string) (*role.Role, error) {
	var role role.Role
	if err := rr.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
