package media

import (
	"bytes"
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	args := m.Called(ctx, file, header, folder)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Delete(ctx context.Context, url string) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

type MockFile struct {
	*bytes.Reader
}

func (f *MockFile) Close() error {
	return nil
}
