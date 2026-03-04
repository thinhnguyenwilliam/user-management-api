package storage

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client *s3.Client
	Bucket string
}

func NewS3Storage(ctx context.Context, region, bucket string) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Storage{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (s *S3Storage) SaveToS3(
	ctx context.Context,
	file multipart.File,
	filename string,
	contentType string,
) error {

	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	return err
}
