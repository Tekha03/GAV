package role

import "context"

type RoleRepository interface {
	GetByName(ctx context.Context, name string) (uint, error)
}
