package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) Storage {
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	extension := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	fullPath := filepath.Join(s.basePath, folder, filename)

	os.MkdirAll(filepath.Join(s.basePath, folder), 0755)

	destination, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		return "", err
	}

	return "/uploads/" + folder + "/" + filename, nil
}

func (s *LocalStorage) Delete(ctx context.Context, url string) error {
	return nil
}
