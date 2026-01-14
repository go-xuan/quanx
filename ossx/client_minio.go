package ossx

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClientBuilder minio客户端构建器
func MinioClientBuilder(config *Config) (Client, error) {
	return NewMinioClient(config)
}

// NewMinioClient 创建minio客户端
func NewMinioClient(config *Config) (*MinioClient, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyId, config.AccessKeySecret, config.AccessToken),
		Secure: false,
		Region: "cn-north-1",
	})
	if err != nil {
		return nil, errorx.Wrap(err, "create minio oss client failed")
	}
	return &MinioClient{config: config, client: client}, nil
}

// 可选项转换
func convertOptions[T any](options ...interface{}) T {
	if len(options) > 0 && options[0] != nil {
		if opts, ok := options[0].(T); ok {
			return opts
		}
	}
	var t T
	return t
}

type MinioClient struct {
	config *Config       // oss配置
	client *minio.Client // minio客户端
}

func (c *MinioClient) GetClient() *minio.Client {
	return c.client
}

func (c *MinioClient) GetInstance() interface{} {
	return c.client
}

func (c *MinioClient) GetConfig() *Config {
	return c.config
}

func (c *MinioClient) Close() error {
	// minio本质上是http客户端，无需关闭
	return nil
}

func (c *MinioClient) CreateBucket(ctx context.Context, bucket string, options ...interface{}) error {
	if exist, err := c.GetClient().BucketExists(ctx, bucket); err != nil {
		return errorx.Wrap(err, "check minio bucket exists failed")
	} else if !exist {
		opts := convertOptions[minio.MakeBucketOptions](options...)
		if err = c.GetClient().MakeBucket(ctx, bucket, opts); err != nil {
			return errorx.Wrap(err, "create minio bucket failed")
		}
	}
	return nil
}

func (c *MinioClient) Upload(ctx context.Context, key string, reader io.Reader, options ...interface{}) error {
	opts := convertOptions[minio.PutObjectOptions](options...)
	info, err := c.GetClient().PutObject(ctx, c.config.Bucket, key, reader, -1, opts)
	if err != nil {
		return errorx.Wrap(err, "upload minio object failed")
	} else if info.Size == 0 {
		return errorx.New("minio upload object size is zero")
	}
	return nil
}

func (c *MinioClient) Get(ctx context.Context, key string, options ...interface{}) (io.ReadCloser, error) {
	opts := convertOptions[minio.GetObjectOptions](options...)
	obj, err := c.GetClient().GetObject(ctx, c.config.Bucket, key, opts)
	if err != nil {
		return nil, errorx.Wrap(err, "get minio object failed")
	}
	return obj, nil
}

func (c *MinioClient) Exist(ctx context.Context, key string, options ...interface{}) (bool, error) {
	opts := convertOptions[minio.StatObjectOptions](options...)
	info, err := c.GetClient().StatObject(ctx, c.config.Bucket, key, opts)
	if err != nil {
		return false, errorx.Wrap(err, "stat minio object failed")
	} else if info.Err != nil {
		return false, errorx.Wrap(info.Err, "minio object not exist")
	}
	return true, nil
}

func (c *MinioClient) Remove(ctx context.Context, key string, options ...interface{}) error {
	opts := convertOptions[minio.RemoveObjectOptions](options...)
	err := c.GetClient().RemoveObject(ctx, c.config.Bucket, key, opts)
	if err != nil {
		return errorx.Wrap(err, "remove minio object failed")
	}
	return nil
}

func (c *MinioClient) GetUrl(ctx context.Context, key string, expires time.Duration, options ...interface{}) (string, error) {
	params := convertOptions[url.Values](options...)
	URL, err := c.GetClient().PresignedGetObject(ctx, c.config.Bucket, key, expires, params)
	if err != nil {
		return "", errorx.Wrap(err, "get minio object url failed")
	}
	return URL.String(), nil
}
