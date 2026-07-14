package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"rag/internal/dto"
	"rag/internal/models"
	"rag/internal/processors"
	"rag/internal/repository"
	"rag/internal/storage"
	"strings"
	"sync"
)

type UploadService struct {
	repo         repository.DatabaseRepo
	storage      storage.FileStorage
	ch           chan *models.UploadedFile
	wg           sync.WaitGroup
	csvProcessor *processors.CSVProcessor
}

func NewUploadService(r repository.DatabaseRepo, s storage.FileStorage, chunkSize int) *UploadService {
	services := &UploadService{
		repo:         r,
		storage:      s,
		ch:           make(chan *models.UploadedFile),
		csvProcessor: processors.NewCSVProcessor(r, chunkSize),
	}
	services.StartWorker()
	return services
}

func (s *UploadService) StartWorker() {
	s.wg.Add(1)
	//spin go rutine
	go s.processFileWorker()
}

// processFileWorker processs uploaded file
func (s *UploadService) processFileWorker() {
	defer s.wg.Done()
	// wait for uploaded files to be added to the channel.
	for uploadedFile := range s.ch {
		if uploadedFile == nil {
			continue
		}
		if err := s.ProcessFile(context.Background(), uploadedFile); err != nil {
			fmt.Println("failed to process file:", err)
		}
	}
}

func (s *UploadService) ProcessFile(ctx context.Context, uploadedFile *models.UploadedFile) error {
	file, err := os.Open(uploadedFile.FilePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	ext := strings.ToLower(filepath.Ext(uploadedFile.FilePath))

	switch ext {
	case ".csv":
		// Mark file as processing
		uploadedFile.Status = "processing"
		uploadedFile.ErrorMessage = ""
		if err := s.csvProcessor.Process(ctx, file, uploadedFile); err != nil {
			// Mark file as failed on db
			uploadedFile.Status = "failed"
			uploadedFile.ErrorMessage = err.Error()
			if err := s.repo.UpdateFile(uploadedFile); err != nil {
				return fmt.Errorf("failed to update uploaded file status: %w", err)
			}

			return fmt.Errorf("failed to process csv file: %w", err)
		}
		// Mark file as proccessed
		uploadedFile.Status = "processed"
		uploadedFile.ErrorMessage = ""
		if err := s.repo.UpdateFile(uploadedFile); err != nil {
			return fmt.Errorf("failed to update uploaded file status: %w", err)
		}

	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
	return nil
}

func (s *UploadService) SaveUploadedFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, description string) (*storage.StoredFile, error) {
	// Check if file already exists before saving it locally.
	_, found, err := s.repo.GetFileByName(header.Filename)
	if err != nil {
		return nil, err
	}
	if found {
		return nil, fmt.Errorf("file already uploaded")
	}
	//save file
	storedFile, err := s.storage.Save(ctx, file, header)
	if err != nil {
		return nil, err
	}
	// save failed on db
	saveFiled, err := s.repo.SaveFileMetaData(storedFile)
	if err != nil {
		return nil, err
	}
	// send saved-file to chan
	s.ch <- saveFiled

	return storedFile, nil
}

func (s *UploadService) StopWorker() {
	close(s.ch)
	s.wg.Wait()
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
