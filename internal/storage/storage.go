package storage

import (
	"context"
	"mime/multipart"
)

type FileStorage interface {
	Save(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*StoredFile, error)
}

type StoredFile struct {
	OriginalFileName string
	StoredFileName   string
	FilePath         string
	Size             int64
	ContentType      string
	Description      string
}
