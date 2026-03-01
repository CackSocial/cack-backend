package storage

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

const maxUploadSize = 10 << 20 // 10MB

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var (
	ErrFileTooLarge      = errors.New("file size exceeds 10MB limit")
	ErrInvalidFileType   = errors.New("invalid file type: only JPEG, PNG, GIF and WebP are allowed")
)

type localStorage struct {
	uploadPath string
	baseURL    string
}

func NewLocalStorage(uploadPath, baseURL string) Storage {
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}
	return &localStorage{uploadPath: uploadPath, baseURL: baseURL}
}

func (s *localStorage) Upload(file *multipart.FileHeader) (string, error) {
	if file.Size > maxUploadSize {
		return "", ErrFileTooLarge
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := src.Close(); err != nil {
			_ = err // read-only source; close error is non-critical
		}
	}()

	// Detect MIME type from file content.
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}
	mimeType := http.DetectContentType(buf[:n])
	if !allowedMIMETypes[mimeType] {
		return "", ErrInvalidFileType
	}
	// Reset reader position after sniffing.
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to seek file: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	dstPath := filepath.Join(s.uploadPath, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if err := dst.Close(); err != nil {
			_ = err // close error logged by caller via io.Copy result
		}
	}()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return s.baseURL + "/uploads/" + filename, nil
}

func (s *localStorage) Delete(filePath string) error {
	fullPath := filepath.Join(s.uploadPath, path.Base(filePath))
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
