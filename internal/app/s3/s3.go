package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/RIBorisov/GophKeeper/internal/config"
	"github.com/RIBorisov/GophKeeper/internal/log"
)

type Client struct {
	client     *minio.Client
	bucketName string
}

func NewS3Client(ctx context.Context, cfg *config.Config) (*Client, error) {
	endpoint := "localhost:9000"
	accessKeyID := "admin"
	secretAccessKey := "password"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start new s3 server: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.S3.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		log.Debug("Bucket not found, going to create it..")
		err = client.MakeBucket(ctx, cfg.S3.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create new bucket: %w", err)
		}
		log.Debug("Successfully created", "bucket", cfg.S3.BucketName)
	}

	return &Client{client: client, bucketName: cfg.S3.BucketName}, nil
}

func (c *Client) PutObject(ctx context.Context, name string, obj io.Reader, size int64) error {
	if _, err := c.client.PutObject(ctx, c.bucketName, name, obj, size, minio.PutObjectOptions{}); err != nil {
		return fmt.Errorf("failed to put object into bucket: %w", err)
	}

	return nil
}

func (c *Client) GetObject(ctx context.Context, name string) ([]byte, error) {
	obj, err := c.client.GetObject(ctx, c.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	defer obj.Close()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(obj); err != nil {
		return nil, fmt.Errorf("failed to read obj from S3: %w", err)
	}
	return buf.Bytes(), nil
}