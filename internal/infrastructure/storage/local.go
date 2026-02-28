package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type localStorage struct {
	uploadPath string
}

func NewLocalStorage(uploadPath string) Storage {
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}
	return &localStorage{uploadPath: uploadPath}
}

func (s *localStorage) Upload(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	dstPath := filepath.Join(s.uploadPath, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return "/uploads/" + filename, nil
}

func (s *localStorage) Delete(filePath string) error {
	fullPath := filepath.Join(s.uploadPath, filepath.Base(filePath))
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
