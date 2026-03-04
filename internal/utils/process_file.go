// user-management-api/internal/utils/process_file.go
// user-management-api/internal/utils/process_file.go
package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/thinhnguyenwilliam/user-management-api/internal/storage"
)

const maxFileSize = 5 << 20 // 5MB

var allowedMimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

func ValidateImage(fileHeader *multipart.FileHeader) (string, error) {

	if fileHeader.Size > maxFileSize {
		return "", errors.New("file too large (max 5MB)")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", errors.New("cannot open file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", errors.New("cannot read file")
	}

	mimeType := http.DetectContentType(buffer[:n])

	ext, allowed := allowedMimeTypes[mimeType]
	if !allowed {
		return "", fmt.Errorf("unsupported file type: %s", mimeType)
	}

	return ext, nil
}

func GenerateFileName(ext string) string {
	return uuid.New().String() + ext
}

func SaveFile(fileHeader *multipart.FileHeader, destination string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func ProcessImageUpload(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	s3Storage *storage.S3Storage,
) (string, error) {

	ext, err := ValidateImage(fileHeader)
	if err != nil {
		return "", err
	}

	filename := GenerateFileName(ext)

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	contentType := mime.TypeByExtension(ext)

	if err := s3Storage.SaveToS3(ctx, file, filename, contentType); err != nil {
		return "", err
	}

	return filename, nil
}
