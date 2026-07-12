package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalStorage struct {
	UploadDir string
}

func NewLocalStorage(uploadDir string) *LocalStorage {
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	return &LocalStorage{
		UploadDir: uploadDir,
	}
}

func (s *LocalStorage) Save(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*StoredFile, error) {
	fmt.Println("************* Saving the file in Weaviate *************")
	if err := os.MkdirAll(s.UploadDir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create upload directory: %w", err)
	}

	originalFileName := filepath.Base(header.Filename)
	storedFileName := buildStoredFileName(originalFileName)
	filePath := filepath.Join(s.UploadDir, storedFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create local file: %w", err)
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("unable to save uploaded file: %w", err)
	}

	return &StoredFile{
		OriginalFileName: originalFileName,
		StoredFileName:   storedFileName,
		FilePath:         filePath,
		Size:             size,
		ContentType:      header.Header.Get("Content-Type"),
	}, nil
}

func buildStoredFileName(original string) string {
	ext := filepath.Ext(original)
	name := strings.TrimSuffix(original, ext)

	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)

	timestamp := time.Now().UnixNano()

	return fmt.Sprintf("%d_%s%s", timestamp, name, ext)
}
