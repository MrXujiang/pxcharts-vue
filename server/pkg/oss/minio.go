package oss

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mime/multipart"
	"mvtable/internal/pkg/config"
)

type Minio struct {
	client *minio.Client
	conf   config.MinioConf
}

func InitMinioClient(c *config.OSS) OSS {
	minioClient, _ := minio.New(c.MinioConf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.MinioConf.AccessKeyId, c.MinioConf.SecretAccessKey, ""),
		Secure: c.MinioConf.UseSSL,
	})

	return &Minio{
		client: minioClient,
		conf:   c.MinioConf,
	}
}

func (c *Minio) UploadFile(file *multipart.FileHeader, directory string) error {
	fileObj, err := file.Open()
	defer fileObj.Close()

	_, err = c.client.PutObject(context.TODO(), c.conf.BucketName, directory+"/"+file.Filename, fileObj, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
	if err != nil {
		return err
	}
	return nil
}

func (c *Minio) DownloadFile(objectName string) (*[]byte, error) {
	obj, err := c.client.GetObject(context.TODO(), c.conf.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	stream, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}
	return &stream, nil
}

func (c *Minio) DeleteFile(objectName string) error {
	err := c.client.RemoveObject(context.TODO(), c.conf.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
