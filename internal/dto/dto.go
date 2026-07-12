package dto

type UploadFileResponse struct {
	UploadedFiles []UploadedFileDTO
}

type UploadedFileDTO struct {
	ID               string
	OriginalFileName string
	StoredFileName   string
	FilePath         string
	Description      string
	ContentType      string
	Status           string
	Size             int64
	CreatedAt        string
	UpdatedAt        string
	ErrorMessage     string
}
