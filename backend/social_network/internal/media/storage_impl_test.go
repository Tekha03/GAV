package media

import (
	"bytes"
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalStorageUploadAndDelete(t *testing.T) {
	basePath := t.TempDir()
	storage := NewLocalStorage(basePath)
	content := []byte("image contents")
	file := &MockFile{bytes.NewReader(content)}
	header := &multipart.FileHeader{Filename: "photo.JPG", Size: int64(len(content))}

	url, err := storage.Upload(context.Background(), file, header, "avatars/user")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(url, "/uploads/avatars/user/") || !strings.HasSuffix(url, ".JPG") {
		t.Fatalf("Upload() url = %q", url)
	}
	path := filepath.Join(basePath, strings.TrimPrefix(url, "/uploads/"))
	stored, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(stored, content) {
		t.Fatalf("stored file = %q, error=%v", stored, err)
	}
	if err := storage.Delete(context.Background(), url); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("deleted file still exists: %v", err)
	}
	if err := storage.Delete(context.Background(), url); err != nil {
		t.Fatalf("deleting missing file error = %v", err)
	}
}

func TestLocalStorageRejectsInvalidDeleteURL(t *testing.T) {
	storage := NewLocalStorage(t.TempDir())
	for _, url := range []string{"", "/uploads/", "/other/file.jpg"} {
		if err := storage.Delete(context.Background(), url); err != ErrInvalidURL {
			t.Fatalf("Delete(%q) error = %v", url, err)
		}
	}
}
