package media

import (
	"context"
	"mime/multipart"
)

type Storage interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (url string, err error)
	Delete(ctx context.Context, url string) error
}
