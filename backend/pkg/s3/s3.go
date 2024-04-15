package s3

import (
	"context"
	"fmt"
	"io"
	"single-window/config"
	"single-window/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const bucketName = "media"

type IS3Client interface {
	PutObjectToBucket(objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (*minio.UploadInfo, error)
}

type S3Client struct {
	S3     *minio.Client
	logger *logger.Logger
	ctx    context.Context
}

var _ IS3Client = (*S3Client)(nil)

func New(ctx context.Context, cfg *config.MINIO, logger *logger.Logger) (*S3Client, error) {
	minioClient, err := minio.New(cfg.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client, err: %w", err)
	}

	s3Client := &S3Client{
		S3:     minioClient,
		logger: logger,
		ctx:    ctx,
	}

	err = s3Client.makeBucket(ctx, minio.MakeBucketOptions{})
	if err != nil {
		return nil, fmt.Errorf(`failed to create bucket with name %s, err: %w`, bucketName, err)
	}

	return s3Client, nil
}

func (s *S3Client) makeBucket(ctx context.Context, bucketOptions minio.MakeBucketOptions) error {
	err := s.S3.MakeBucket(ctx, bucketName, bucketOptions)
	if err != nil {
		exists, errBucketExists := s.S3.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			s.logger.Info(`bucket with name "%s" is already exists`, bucketName)
			return nil
		} else {
			return fmt.Errorf("failed to check if bucket exists, err: %w", err)
		}
	}

	return nil
}

func (s *S3Client) PutObjectToBucket(objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (*minio.UploadInfo, error) {
	info, err := s.S3.PutObject(s.ctx, bucketName, objectName, reader, objectSize, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to put object in s3, err: %w", err)
	}

	s.logger.Info("successfully put object to s3, info: %+v", info)

	return &info, nil
}
