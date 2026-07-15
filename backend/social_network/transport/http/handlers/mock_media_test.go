package handlers

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MockMediaService struct {
	mock.Mock
}

func (m *MockMediaService) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, path string) (string, error) {
	args := m.Called(ctx, file, header, path)
	return args.String(0), args.Error(1)
}

func (m *MockMediaService) Delete(ctx context.Context, url string) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}
