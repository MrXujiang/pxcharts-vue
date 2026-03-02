package oss

import (
	"fmt"
	"mime/multipart"
	"mvtable/internal/pkg/config"
	"sync"
)

type OSS interface {
	UploadFile(file *multipart.FileHeader, directory string) error
	DownloadFile(objectName string) (*[]byte, error)
	DeleteFile(objectName string) error
}

type Client struct {
	MinioConf  config.MinioConf
	AliyunConf config.AliyunConf
	Provider   string
	Client     OSS
	Aliyun     OSS
	Minio      OSS
}

var (
	client *Client
	once   sync.Once
)

func Init(c *config.OSS) (err error) {
	once.Do(func() {
		switch c.Provider {
		case "minio":
			client = &Client{
				MinioConf:  c.MinioConf,
				AliyunConf: c.AliyunConf,
				Provider:   c.Provider,
				Client:     InitMinioClient(c),
				Aliyun:     InitAliyunClient(c),
				Minio:      InitMinioClient(c),
			}
			return
		case "aliyun":
			client = &Client{
				MinioConf:  c.MinioConf,
				AliyunConf: c.AliyunConf,
				Provider:   c.Provider,
				Client:     InitAliyunClient(c),
				Aliyun:     InitAliyunClient(c),
				Minio:      InitMinioClient(c),
			}
			return

		default:
			err = fmt.Errorf("unsupported provider: %s", c.Provider)
			return
		}
	})
	return err
}

func GetClient() *Client {
	return client
}

func (s *Client) UploadFile(file *multipart.FileHeader, directory string) error {
	return s.Client.UploadFile(file, directory)
}

func (s *Client) DownloadFile(objectName string) (*[]byte, error) {
	return s.Client.DownloadFile(objectName)
}

func (s *Client) DeleteFile(objectName string) error {
	return s.Client.DeleteFile(objectName)
}
