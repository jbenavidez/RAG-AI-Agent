package services

import (
	"context"
	"mime/multipart"
	"rag/internal/dto"
	dbrepo "rag/internal/repository/db_repo"
	"rag/internal/storage"
)

type UploadService struct {
	repo    *dbrepo.WeaviateDBRepo
	storage storage.FileStorage
}

func NewUploadService(r *dbrepo.WeaviateDBRepo, s storage.FileStorage) *UploadService {
	return &UploadService{
		repo:    r,
		storage: s,
	}
}

func (s *UploadService) SaveUploadedFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, description string) (*storage.StoredFile, error) {
	//save file
	storedFile, err := s.storage.Save(ctx, file, header)
	if err != nil {
		return nil, err
	}
	// save failed on db
	err = s.repo.SaveFileMetaData(storedFile)
	if err != nil {
		return nil, err
	}

	return storedFile, nil
}

func (s *UploadService) GetAllUploadedFile() (*dto.UploadFileResponse, error) {
	// get all uploaded files
	files, err := s.repo.GetAllUploadFiles()
	if err != nil {
		return nil, err
	}
	/// models to dtos
	uploadedFiles := make([]dto.UploadedFileDTO, len(files))
	for i := range files {
		file := files[i]
		uploadedFiles[i] = dto.UploadedFileDTO{
			ID:               file.ID,
			OriginalFileName: file.OriginalFileName,
			StoredFileName:   file.StoredFileName,
			FilePath:         file.FilePath,
			Description:      file.Description,
			ContentType:      file.ContentType,
			Status:           file.Status,
			Size:             file.Size,
			CreatedAt:        file.CreatedAt.Format("2006-01-02 15:04"),
			UpdatedAt:        file.UpdatedAt.Format("2006-01-02 15:04"),
			ErrorMessage:     file.ErrorMessage,
		}
	}
	//
	resp := dto.UploadFileResponse{
		UploadedFiles: uploadedFiles,
	}

	return &resp, nil
}
