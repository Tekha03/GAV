package repository

import "context"

type TransactionManager interface {
	WithPrivateChatLock(ctx context.Context, privateKey string, fn func(context.Context) error) error
}
