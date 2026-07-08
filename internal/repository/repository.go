package repository

import "github.com/weaviate/weaviate-go-client/v5/weaviate"

type DatabaseRepo interface {
	Connection() *weaviate.Client
}
