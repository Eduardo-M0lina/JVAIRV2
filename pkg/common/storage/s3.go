package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appconfig "github.com/angumol/jvairv2/configs"
)

// S3Client maneja las operaciones de almacenamiento en S3
type S3Client struct {
	client *s3.Client
	bucket string
	url    string // CDN URL override (opcional)
}

// NewS3Client crea un nuevo cliente S3
func NewS3Client(cfg *appconfig.S3Config) (*S3Client, error) {
	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
	}

	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var s3Opts []func(*s3.Options)
	if cfg.Endpoint != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		})
	}

	client := s3.NewFromConfig(awsCfg, s3Opts...)

	return &S3Client{
		client: client,
		bucket: cfg.Bucket,
		url:    cfg.URL,
	}, nil
}

// Upload sube un archivo a S3
func (s *S3Client) Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to upload file to S3",
			slog.String("error", err.Error()),
			slog.String("key", key))
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	var fileURL string
	if s.url != "" {
		fileURL = fmt.Sprintf("%s/%s", strings.TrimRight(s.url, "/"), key)
	} else {
		fileURL = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
	}
	return fileURL, nil
}

// Delete elimina un archivo de S3
func (s *S3Client) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete file from S3",
			slog.String("error", err.Error()),
			slog.String("key", key))
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetObject obtiene un archivo de S3
func (s *S3Client) GetObject(ctx context.Context, key string) (io.ReadCloser, string, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get file from S3",
			slog.String("error", err.Error()),
			slog.String("key", key))
		return nil, "", fmt.Errorf("failed to get file: %w", err)
	}

	contentType := ""
	if result.ContentType != nil {
		contentType = *result.ContentType
	}

	return result.Body, contentType, nil
}

// GetBucket retorna el nombre del bucket
func (s *S3Client) GetBucket() string {
	return s.bucket
}
