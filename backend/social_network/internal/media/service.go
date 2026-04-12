package media

import (
	"context"
	"mime/multipart"
)

type MediaService interface {
	UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error)
	Delete(ctx context.Context, url string) error
}
