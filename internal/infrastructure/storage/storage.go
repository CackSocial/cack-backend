package storage

import "mime/multipart"

type Storage interface {
	Upload(file *multipart.FileHeader) (string, error)
	Delete(filePath string) error
}
