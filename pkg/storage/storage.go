package storage

import (
	"context"
	"io"
)

type Storage interface {
	UploadObject(ctx context.Context, name string, data io.Reader) (string, error)
}
