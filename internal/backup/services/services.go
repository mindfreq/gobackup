package services

import (
	"io"
	"context"
)


type Service interface {
	Upload(ctx context.Context, key string, body io.Reader) error
}