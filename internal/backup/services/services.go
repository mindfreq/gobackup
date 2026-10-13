package services

import (
	"context"
	"io"
)

type Service interface {
	Upload(ctx context.Context, key string, body io.Reader) error
}
