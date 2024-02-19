// Package fileutil ...
package fileutil

import (
	"io"
	"net/http"
)

// GetMimeTypeFile ...
func GetMimeTypeFile(data io.ReadSeeker) (string, error) {
	fileHeader := make([]byte, 512)
	if _, err := data.Read(fileHeader); err != nil {
		return "", err
	}

	// Set position back to start.
	if _, err := data.Seek(0, 0); err != nil {
		return "", err
	}

	mime := http.DetectContentType(fileHeader)

	return mime, nil
}
