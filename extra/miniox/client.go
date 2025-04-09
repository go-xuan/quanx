package miniox

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
)

const Region = "cn-north-1"

var _client *Client

func this() *Client {
	if _client == nil {
		panic("minio client not initialized, please check the relevant config")
	}
	return _client
}

type Client struct {
	config *Config
	client *minio.Client // minio客户端
}

func (c *Client) Instance() *minio.Client {
	return c.client
}

func (c *Client) Config() *Config {
	return c.config
}

// CreateBucket 创建桶
func (c *Client) CreateBucket(ctx context.Context, name string) error {
	var client = c.Instance()
	if exist, err := client.BucketExists(ctx, name); err != nil {
		return errorx.Wrap(err, "check bucket exists error")
	} else if !exist {
		if err = client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: Region}); err != nil {
			return errorx.Wrap(err, "make bucket error")
		}
		if err = client.SetBucketPolicy(ctx, name, defaultBucketPolicy(name)); err != nil {
			return errorx.Wrap(err, "set bucket policy error")
		}
	}
	return nil
}

// PutObject 上传文件
func (c *Client) PutObject(ctx context.Context, bucketName, minioPath string, reader io.Reader) error {
	if _, err := c.Instance().PutObject(ctx, bucketName, minioPath, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"}); err != nil {
		return errorx.Wrap(err, "put object error")
	}
	return nil
}

// ObjectExist 获取对象是否存在
func (c *Client) ObjectExist(ctx context.Context, bucketName string, minioPath string) (bool, error) {
	if objInfo, err := c.Instance().StatObject(ctx, bucketName, minioPath, minio.StatObjectOptions{}); err != nil {
		return false, errorx.Wrap(err, "stat object error")
	} else {
		return objInfo.Size > 0, nil
	}
}

// UploadFile 上传文件
func (c *Client) UploadFile(ctx context.Context, bucketName string, minioPath string, file *multipart.FileHeader) error {
	if exist, err := c.ObjectExist(ctx, c.Config().BucketName, minioPath); err != nil {
		return errorx.Wrap(err, "check object exist error")
	} else if !exist {
		var f multipart.File
		if f, err = file.Open(); err != nil {
			return errorx.Wrap(err, "open file error")
		}
		defer f.Close()
		if err = c.PutObject(ctx, bucketName, minioPath, f); err != nil {
			return errorx.Wrap(err, "put object error")
		}
	}
	return nil
}

// FGetObject 下载文件
func (c *Client) FGetObject(ctx context.Context, bucketName, minioPath, savePath string) error {
	if err := c.Instance().FGetObject(ctx, bucketName, minioPath, savePath, minio.GetObjectOptions{}); err != nil {
		return errorx.Wrap(err, "get object error")
	}
	return nil
}

// RemoveObject 删除文件
func (c *Client) RemoveObject(ctx context.Context, bucketName, minioPath string) error {
	if err := c.Instance().RemoveObject(ctx, bucketName, minioPath, minio.RemoveObjectOptions{GovernanceBypass: true}); err != nil {
		return errorx.Wrap(err, "remove object error")
	}
	return nil
}

// RemoveObjectBatch 通过文件名称删除文件
func (c *Client) RemoveObjectBatch(ctx context.Context, bucketName string, minioPaths []string) error {
	for _, minioPath := range minioPaths {
		if err := c.RemoveObject(ctx, bucketName, minioPath); err != nil {
			return errorx.Wrap(err, "remove object error")
		}
	}
	return nil
}

// PresignedGetObject 下载链接
func (c *Client) PresignedGetObject(ctx context.Context, minioPath string) (string, error) {
	var expires = time.Duration(c.Config().Expire) * time.Minute
	var bucketName = c.Config().BucketName
	if URL, err := c.Instance().PresignedGetObject(ctx, bucketName, minioPath, expires, nil); err != nil {
		return "", errorx.Wrap(err, "presigned get object error")
	} else {
		return URL.String(), nil
	}
}

// PresignedGetObjects 下载链接
func (c *Client) PresignedGetObjects(ctx context.Context, minioPaths []string) ([]string, error) {
	var urls = make([]string, 0)
	for _, minioPath := range minioPaths {
		if object, err := c.PresignedGetObject(ctx, minioPath); err != nil {
			return nil, errorx.Wrap(err, "presigned get object error")
		} else {
			urls = append(urls, object)
		}
	}
	return urls, nil
}

// UploadFileByUrl 通过文件路径上传文件到桶
func (c *Client) UploadFileByUrl(ctx context.Context, bucketName string, fileName string, url string) (string, error) {
	if fileBytes, err := filex.GetFileBytesByUrl(url); err != nil {
		return "", errorx.Wrap(err, "get file bytes error")
	} else {
		minioPath := c.Config().MinioPath(fileName)
		if err = c.PutObject(ctx, bucketName, minioPath, bytes.NewBuffer(fileBytes)); err != nil {
			return minioPath, errorx.Wrap(err, "put object error")
		}
		return minioPath, nil
	}
}
