package oss

import (
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"io"
	"mime/multipart"
	"mvtable/internal/pkg/config"
)

type Aliyun struct {
	client *oss.Client
	conf   config.AliyunConf
}

func InitAliyunClient(c *config.OSS) OSS {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AliyunConf.AccessKeyID, c.AliyunConf.AccessKeySecret)).
		WithRegion(c.AliyunConf.Region)

	return &Aliyun{
		client: oss.NewClient(cfg),
		conf:   c.AliyunConf,
	}
}

// UploadFile 上传文件到 OSS
func (c *Aliyun) UploadFile(file *multipart.FileHeader, directory string) error {
	fileObj, err := file.Open()
	defer fileObj.Close()

	contentType := file.Header.Get("Content-Type")
	request := &oss.PutObjectRequest{
		Bucket:       oss.Ptr(c.conf.BucketName),
		Key:          oss.Ptr(directory + "/" + file.Filename),
		Body:         fileObj,
		StorageClass: oss.StorageClassStandard,
		Acl:          oss.ObjectACLPublicRead,
		ContentType:  &contentType,
	}

	result, err := c.client.PutObject(context.TODO(), request)
	if err != nil {
		return err
	}

	if result.Status != "200 OK" {
		return fmt.Errorf("upload fail: %#v\n", result)
	}

	return nil
}

// DownloadFile 获取文件
func (c *Aliyun) DownloadFile(objectName string) (*[]byte, error) {
	request := &oss.GetObjectRequest{
		Bucket: oss.Ptr(c.conf.BucketName),
		Key:    oss.Ptr(objectName),
	}

	result, err := c.client.GetObject(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	data, _ := io.ReadAll(result.Body)
	return &data, nil
}

func (c *Aliyun) DeleteFile(objectName string) error {
	request := &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(c.conf.BucketName),
		Key:    oss.Ptr(objectName),
	}
	result, err := c.client.DeleteObject(context.TODO(), request)
	if err != nil {
		return err
	}

	if !(result.Status == "204 No Content" || result.Status == "200 OK") {
		return fmt.Errorf("delete fail: %v\n", result)
	}
	return nil
}
