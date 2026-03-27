package media

import (
	"context"
	"mime/multipart"
	"path/filepath"
)

const FIVE_MEGABYTES int64 = 5 * 1024 * 1024

type service struct {
	storage Storage
}

func NewService(storage Storage) (MediaService, error) {
	if storage == nil {
		return nil, ErrStorageNil
	}
	return &service{storage: storage}, nil
}

func (s *service) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	if header.Size > FIVE_MEGABYTES {
		return "", ErrFileTooLarge
	}

	extension := filepath.Ext(header.Filename)
	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" && extension != ".webp" {
		return "", ErrInvalidFileType
	}

	return s.storage.Upload(ctx, file, header, folder)
}


