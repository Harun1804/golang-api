package minio

import (
	"context"
	"fmt"
	"galaxy/backend-api/config"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioHelper struct {
  Client *minio.Client
  BucketName string
  Ctx context.Context
}

func InitMinio() (*MinioHelper, error) {
  endpoint := config.GetEnv("MINIO_ENDPOINT", "localhost:9000")
  accessKey := config.GetEnv("MINIO_ACCESS_KEY", "minioadmin")
  secretKey := config.GetEnv("MINIO_SECRET_KEY", "minioadmin")
  useSSL := false
  bucketName := config.GetEnv("MINIO_BUCKET", "latihan-bucket")
  ctx := context.Background()

  client, err := minio.New(
  endpoint,
  &minio.Options{
      Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
      Secure: useSSL,
    },
  )
  if err != nil {
    return nil, err
  }

  fmt.Println("MinIO connected successfully!")

  helper := &MinioHelper{
    Client: client,
    BucketName: bucketName,
    Ctx: ctx,
  }

  if err := helper.CreateAndCheckBucket(); err != nil {
    return nil, err
  }

  return helper, nil
}

func (m *MinioHelper) CreateAndCheckBucket() error {
  location := config.GetEnv("MINIO_REGION", "us-east-1")
  exists, err := m.Client.BucketExists(m.Ctx, m.BucketName)
  if err != nil {
    return err
  }

  if !exists {
    err = m.Client.MakeBucket(m.Ctx, m.BucketName, minio.MakeBucketOptions{Region: location})
    if err != nil {
      return err
    }

    fmt.Printf("Bucket %s created successfully\n", m.BucketName)
  } else {
    fmt.Printf("Bucket %s already exists\n", m.BucketName)
  }

  return nil
}

func (m *MinioHelper) UploadFile(objectName string, file io.Reader) error {
  _, err := m.Client.PutObject(m.Ctx, m.BucketName, objectName, file, -1, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
  return err
}

func (m *MinioHelper) DeleteFile(objectName string) error {
  err := m.Client.RemoveObject(m.Ctx, m.BucketName, objectName, minio.RemoveObjectOptions{})
  if err != nil {
    return fmt.Errorf("failed to delete file %s: %w", objectName, err)
  }
  return nil
}

func (m *MinioHelper) GetFileUrl(objectName string) (string, error) {
    stat, err := m.Client.StatObject(m.Ctx, m.BucketName, objectName, minio.StatObjectOptions{})
    if err != nil {
        return "", err
    }
    reqParams := make(url.Values)
    reqParams.Set("response-content-disposition", "inline")
    if stat.ContentType != "" {
        reqParams.Set("response-content-type", stat.ContentType)
    }

    presignedURL, err := m.Client.PresignedGetObject(m.Ctx, m.BucketName, objectName, time.Second*24*60*60, reqParams)
    if err != nil {
        return "", err
    }
    return presignedURL.String(), nil
}

func (m *MinioHelper) DownloadFile(objectName string, w http.ResponseWriter) (error) {
  object, err := m.Client.GetObject(m.Ctx, m.BucketName, objectName, minio.GetObjectOptions{})
  if err != nil {
    return err
  }
  defer object.Close()
  w.Header().Set("Content-Disposition", "attachment; filename="+objectName)
  w.Header().Set("Content-Type", "application/octet-stream")
  _, err = io.Copy(w, object)
  return err
}