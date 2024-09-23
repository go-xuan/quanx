package miniox

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
)

const Region = "cn-north-1"

var _handler *Handler

func this() *Handler {
	if _handler == nil {
		panic("the minio handler has not been initialized, please check the relevant config")
	}
	return _handler
}

func GetConfig() *Config {
	return this().config
}

func Client() *minio.Client {
	return this().client
}

// CreateBucket 创建桶
func CreateBucket(ctx context.Context, name string) error {
	var client = Client()
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
func PutObject(ctx context.Context, bucketName, minioPath string, reader io.Reader) error {
	var client = Client()
	if _, err := client.PutObject(ctx, bucketName, minioPath, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"}); err != nil {
		return errorx.Wrap(err, "put object error")
	}
	return nil
}

// ObjectExist 获取对象是否存在
func ObjectExist(ctx context.Context, bucketName string, minioPath string) (bool, error) {
	if objInfo, err := Client().StatObject(ctx, bucketName, minioPath, minio.StatObjectOptions{}); err != nil {
		return false, errorx.Wrap(err, "stat object error")
	} else {
		return objInfo.Size > 0, nil
	}
}

// UploadFile 上传文件
func UploadFile(ctx context.Context, bucketName string, minioPath string, file *multipart.FileHeader) error {
	if exist, err := ObjectExist(ctx, GetConfig().BucketName, minioPath); err != nil {
		return errorx.Wrap(err, "check object exist error")
	} else if !exist {
		var f multipart.File
		if f, err = file.Open(); err != nil {
			return errorx.Wrap(err, "open file error")
		}
		defer f.Close()
		if err = PutObject(ctx, bucketName, minioPath, f); err != nil {
			return errorx.Wrap(err, "put object error")
		}
	}
	return nil
}

// FGetObject 下载文件
func FGetObject(ctx context.Context, bucketName, minioPath, savePath string) error {
	var client = Client()
	if err := client.FGetObject(ctx, bucketName, minioPath, savePath, minio.GetObjectOptions{}); err != nil {
		return errorx.Wrap(err, "get object error")
	}
	return nil
}

// RemoveObject 删除文件
func RemoveObject(ctx context.Context, bucketName, minioPath string) error {
	var client = Client()
	if err := client.RemoveObject(ctx, bucketName, minioPath, minio.RemoveObjectOptions{GovernanceBypass: true}); err != nil {
		return errorx.Wrap(err, "remove object error")
	}
	return nil
}

// RemoveObjectBatch 通过文件名称删除文件
func (h *Handler) RemoveObjectBatch(ctx context.Context, bucketName string, minioPaths []string) error {
	for _, minioPath := range minioPaths {
		if err := RemoveObject(ctx, bucketName, minioPath); err != nil {
			return errorx.Wrap(err, "remove object error")
		}
	}
	return nil
}

// PresignedGetObject 下载链接
func PresignedGetObject(ctx context.Context, minioPath string) (string, error) {
	var expires = time.Duration(GetConfig().Expire) * time.Minute
	var bucketName = GetConfig().BucketName
	if URL, err := Client().PresignedGetObject(ctx, bucketName, minioPath, expires, nil); err != nil {
		return "", errorx.Wrap(err, "presigned get object error")
	} else {
		return URL.String(), nil
	}
}

// PresignedGetObjects 下载链接
func PresignedGetObjects(ctx context.Context, minioPaths []string) ([]string, error) {
	var urls = make([]string, 0)
	for _, minioPath := range minioPaths {
		if object, err := PresignedGetObject(ctx, minioPath); err != nil {
			return nil, errorx.Wrap(err, "presigned get object error")
		} else {
			urls = append(urls, object)
		}
	}
	return urls, nil
}

// UploadFileByUrl 通过文件路径上传文件到桶
func UploadFileByUrl(ctx context.Context, bucketName string, fileName string, url string) (minioPath string, err error) {
	var fileBytes []byte
	if fileBytes, err = filex.GetFileBytesByUrl(url); err != nil {
		return
	}
	minioPath = GetConfig().MinioPath(fileName)
	if err = PutObject(ctx, bucketName, minioPath, bytes.NewBuffer(fileBytes)); err != nil {
		return
	}
	return
}

// Handler minio控制器
type Handler struct {
	config *Config       // minio配置
	client *minio.Client // minio客户端
}

func (h *Handler) GetConfig() *Config {
	return h.config
}

func (h *Handler) GetClient() *minio.Client {
	return h.client
}
