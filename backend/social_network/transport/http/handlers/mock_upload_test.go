package handlers

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MockUploadService struct {
	mock.Mock
}

func (m *MockUploadService) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	args := m.Called(ctx, file, header, folder)
	return args.String(0), args.Error(1)
}

func (m *MockUploadService) Delete(ctx context.Context, url string) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}
