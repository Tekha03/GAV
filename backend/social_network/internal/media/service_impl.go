package media

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"
)

const FIVE_MEGABYTES int64 = 5 * 1024 * 1024
const TWENTY_FIVE_MEGABYTES int64 = 25 * 1024 * 1024

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

	extension := strings.ToLower(filepath.Ext(header.Filename))
	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" && extension != ".webp" {
		return "", ErrInvalidFileType
	}

	return s.storage.Upload(ctx, file, header, folder)
}

func (s *service) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	if header.Size > TWENTY_FIVE_MEGABYTES {
		return "", ErrAttachmentTooLarge
	}

	extension := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]struct{}{
		".jpg": {}, ".jpeg": {}, ".png": {}, ".webp": {}, ".heic": {},
		".mp4": {}, ".mov": {}, ".m4v": {},
		".m4a": {}, ".mp3": {}, ".wav": {}, ".aac": {},
		".pdf": {}, ".txt": {}, ".rtf": {}, ".csv": {},
		".doc": {}, ".docx": {}, ".xls": {}, ".xlsx": {},
		".ppt": {}, ".pptx": {}, ".zip": {},
	}
	if _, ok := allowed[extension]; !ok {
		return "", ErrInvalidFileType
	}

	return s.storage.Upload(ctx, file, header, folder)
}

func (s *service) Delete(ctx context.Context, url string) error {
	if url == "" {
		return ErrEmptyURL
	}

	return s.storage.Delete(ctx, url)
}
