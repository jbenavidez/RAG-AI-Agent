package repository

type DatabaseRepo interface {
	GetTotalDocs() (int, error)
}
