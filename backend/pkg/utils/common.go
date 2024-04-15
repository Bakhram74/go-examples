package utils

import (
	"bytes"
	"fmt"
	"single-window/pkg/s3"

	"github.com/minio/minio-go/v7"
	"github.com/oapi-codegen/runtime/types"
)

func Ptr[T any](v T) *T {
	return &v
}

// TODO: перенести всю логику в s3Client
func UploadFile(s3Client s3.IS3Client, folderName string, file *types.File) (*string, error) {
	var attachmentPath *string

	if file != nil {
		fBytes, err := file.Bytes()
		if err != nil {
			return nil, fmt.Errorf("failed to get file bytes, err: %w", err)
		}

		filePath := fmt.Sprintf("%s/%s", folderName, file.Filename())
		// TODO: перед этим, нужно проверить, нет ли там уже файла с таким же названием
		uploadInfo, err := s3Client.PutObjectToBucket(filePath, bytes.NewReader(fBytes), int64(len(fBytes)), minio.PutObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to put file in bucket, err: %w", err)
		}
		attachmentPath = &uploadInfo.Key
	}

	return attachmentPath, nil
}
