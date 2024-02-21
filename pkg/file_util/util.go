// Package fileutil provides utility functions for working with files.
package fileutil

import (
	"io"
	"net/http"
)

// GetMimeTypeFile detects the MIME type of a file based on its content.
// It takes an io.ReadSeeker as input, reads the first 512 bytes, and uses http.DetectContentType
// to determine the MIME type. It then resets the position of the io.ReadSeeker to the beginning.
// Parameters:
//   - data: The io.ReadSeeker representing the file content.
//
// Returns:
//   - string: The detected MIME type of the file.
//   - error: An error, if any, encountered during the process.
func GetMimeTypeFile(data io.ReadSeeker) (string, error) {
	// Read the first 512 bytes of the file content.
	fileHeader := make([]byte, 512)
	if _, err := data.Read(fileHeader); err != nil {
		return "", err
	}

	// Set the position back to the start.
	if _, err := data.Seek(0, 0); err != nil {
		return "", err
	}

	// Detect the MIME type based on the file content.
	mime := http.DetectContentType(fileHeader)

	return mime, nil
}
