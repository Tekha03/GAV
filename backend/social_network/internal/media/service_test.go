package media

import (
	"bytes"
	"context"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		storage := new(MockStorage)

		s, err := NewService(storage)

		require.NoError(t, err)
		require.NotNil(t, s)
	})

	t.Run("nil storage", func(t *testing.T) {
		s, err := NewService(nil)

		require.ErrorIs(t, err, ErrStorageNil)
		require.Nil(t, s)
	})
}

func TestService_UploadImage(t *testing.T) {
	ctx := context.Background()

	createFile := func(filename string, size int64) (multipart.File, *multipart.FileHeader) {
		data := make([]byte, size)
		reader := bytes.NewReader(data)
		file := &MockFile{reader}

		header := &multipart.FileHeader{
			Filename: filename,
			Size:     size,
		}

		return file, header
	}

	t.Run("file too large", func(t *testing.T) {
		storage := new(MockStorage)
		s, _ := NewService(storage)

		file, header := createFile("test.jpg", FIVE_MEGABYTES+1)

		url, err := s.UploadImage(ctx, file, header, "avatars")

		require.ErrorIs(t, err, ErrFileTooLarge)
		require.Empty(t, url)
	})

	t.Run("invalid file type", func(t *testing.T) {
		storage := new(MockStorage)
		s, _ := NewService(storage)

		file, header := createFile("test.txt", 100)

		url, err := s.UploadImage(ctx, file, header, "avatars")

		require.ErrorIs(t, err, ErrInvalidFileType)
		require.Empty(t, url)
	})

	t.Run("storage error", func(t *testing.T) {
		storage := new(MockStorage)
		s, _ := NewService(storage)

		file, header := createFile("test.jpg", 100)

		storage.
			On("Upload", ctx, file, header, "avatars").
			Return("", ErrStorageNil).
			Once()

		url, err := s.UploadImage(ctx, file, header, "avatars")

		require.Error(t, err)
		require.Empty(t, url)
	})

	t.Run("success jpg", func(t *testing.T) {
		storage := new(MockStorage)
		s, _ := NewService(storage)

		file, header := createFile("test.jpg", 100)

		storage.
			On("Upload", ctx, file, header, "avatars").
			Return("/uploads/avatars/test.jpg", nil).
			Once()

		url, err := s.UploadImage(ctx, file, header, "avatars")

		require.NoError(t, err)
		require.Equal(t, "/uploads/avatars/test.jpg", url)
	})

	t.Run("success png", func(t *testing.T) {
		storage := new(MockStorage)
		s, _ := NewService(storage)

		file, header := createFile("test.png", 100)

		storage.
			On("Upload", ctx, file, header, "avatars").
			Return("/uploads/avatars/test.png", nil).
			Once()

		url, err := s.UploadImage(ctx, file, header, "avatars")

		require.NoError(t, err)
		require.Equal(t, "/uploads/avatars/test.png", url)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("empty url", func(t *testing.T) {
		storage := new(MockStorage)
		s, _ := NewService(storage)

		err := s.Delete(ctx, "")

		require.Error(t, err)
	})

	t.Run("storage error", func(t *testing.T) {
		storage := new(MockStorage)
		s, _ := NewService(storage)

		url := "/uploads/avatars/test.jpg"

		storage.
			On("Delete", ctx, url).
			Return(ErrStorageNil).
			Once()

		err := s.Delete(ctx, url)

		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		storage := new(MockStorage)
		s, _ := NewService(storage)

		url := "/uploads/avatars/test.jpg"

		storage.
			On("Delete", ctx, url).
			Return(nil).
			Once()

		err := s.Delete(ctx, url)

		require.NoError(t, err)
	})
}

func TestService_UploadFile(t *testing.T) {
	ctx := context.Background()
	createFile := func(filename string, size int64) (multipart.File, *multipart.FileHeader) {
		file := &MockFile{bytes.NewReader(make([]byte, min(size, 16)))}
		return file, &multipart.FileHeader{Filename: filename, Size: size}
	}

	t.Run("file too large", func(t *testing.T) {
		service, _ := NewService(new(MockStorage))
		file, header := createFile("video.mp4", TWENTY_FIVE_MEGABYTES+1)
		url, err := service.UploadFile(ctx, file, header, "chat")
		require.ErrorIs(t, err, ErrAttachmentTooLarge)
		require.Empty(t, url)
	})

	t.Run("invalid extension", func(t *testing.T) {
		service, _ := NewService(new(MockStorage))
		file, header := createFile("payload.exe", 10)
		url, err := service.UploadFile(ctx, file, header, "chat")
		require.ErrorIs(t, err, ErrInvalidFileType)
		require.Empty(t, url)
	})

	for _, filename := range []string{"photo.HEIC", "voice.m4a", "video.mov", "document.pdf", "archive.zip"} {
		t.Run("success "+filename, func(t *testing.T) {
			storage := new(MockStorage)
			service, _ := NewService(storage)
			file, header := createFile(filename, 10)
			storage.On("Upload", ctx, file, header, "chat").Return("/uploads/chat/"+filename, nil).Once()
			url, err := service.UploadFile(ctx, file, header, "chat")
			require.NoError(t, err)
			require.Equal(t, "/uploads/chat/"+filename, url)
		})
	}

	t.Run("storage error", func(t *testing.T) {
		storage := new(MockStorage)
		service, _ := NewService(storage)
		file, header := createFile("photo.jpg", 10)
		storage.On("Upload", ctx, file, header, "chat").Return("", ErrStorageNil).Once()
		url, err := service.UploadFile(ctx, file, header, "chat")
		require.ErrorIs(t, err, ErrStorageNil)
		require.Empty(t, url)
	})
}
