package miniox

import (
	"bytes"
	"context"
	"github.com/go-xuan/quanx/utilx/filex"
	"io"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

var handler *Handler

// minio控制器
type Handler struct {
	Config *Minio        // minio配置
	Client *minio.Client // minio客户端
}

func This() *Handler {
	if handler == nil {
		panic("The minio handler has not been initialized, please check the relevant config")
	}
	return handler
}

// 创建桶
func (h *Handler) CreateBucket(ctx context.Context, bucketName string) (err error) {
	var exist bool
	if exist, err = h.Client.BucketExists(ctx, bucketName); err != nil {
		return
	}
	if !exist {
		if err = h.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: ""}); err != nil {
			log.Warnf("create bucket error , %+v", err)
			return
		}
		if err = h.Client.SetBucketPolicy(ctx, bucketName, defaultBucketPolicy(bucketName)); err != nil {
			log.Warnf("set bucket policy error , %+v", err)
			return
		}
	}
	return
}

// 生成minio存储路径
func (h *Handler) NewMinioPath(fileName string) string {
	return h.Config.MinioPath(fileName)
}

// 上传文件
func (h *Handler) PutObject(ctx context.Context, bucketName, minioPath string, reader io.Reader) (err error) {
	_, err = h.Client.PutObject(ctx, bucketName, minioPath, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return
	}
	return
}

// 下载文件
func (h *Handler) FGetObject(ctx context.Context, bucketName, minioPath, savePath string) (err error) {
	err = h.Client.FGetObject(ctx, bucketName, minioPath, savePath, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	return
}

// 删除文件
func (h *Handler) RemoveObject(ctx context.Context, bucketName, minioPath string) (err error) {
	err = h.Client.RemoveObject(ctx, bucketName, minioPath, minio.RemoveObjectOptions{GovernanceBypass: true})
	if err != nil {
		return
	}
	return
}

// 下载链接
func (h *Handler) PresignedGetObject(ctx context.Context, minioPath string) (minioUrl string, err error) {
	var URL *url.URL
	expires := time.Duration(h.Config.Expire) * time.Minute
	URL, err = h.Client.PresignedGetObject(ctx, h.Config.BucketName, minioPath, expires, nil)
	if err != nil {
		return
	}
	minioUrl = URL.String()
	return
}

// 通过文件路径上传文件到桶
func (h *Handler) UploadFileByUrl(ctx context.Context, bucketName string, fileName string, url string) (minioPath string, err error) {
	var fileBytes []byte
	if fileBytes, err = filex.GetFileBytesByUrl(url); err != nil {
		return
	}
	minioPath = h.NewMinioPath(fileName)

	if err = h.PutObject(ctx, bucketName, minioPath, bytes.NewBuffer(fileBytes)); err != nil {
		return
	}
	return
}

// 上传文件
func (h *Handler) UploadFile(ctx context.Context, bucketName string, minioPath string, file *multipart.FileHeader) (err error) {
	var exist bool
	if exist, err = h.ObjectExist(ctx, h.Config.BucketName, minioPath); err != nil {
		return
	}
	if !exist {
		var mf multipart.File
		if mf, err = file.Open(); err != nil {
			return
		}
		defer mf.Close()
		if err = h.PutObject(ctx, bucketName, minioPath, mf); err != nil {
			return
		}
	}
	return
}

// 获取对象是否存在
func (h *Handler) ObjectExist(ctx context.Context, bucketName string, minioPath string) (exist bool, err error) {
	var objInfo minio.ObjectInfo
	if objInfo, err = h.Client.StatObject(ctx, bucketName, minioPath, minio.StatObjectOptions{}); err != nil {
		return
	}
	exist = objInfo.Size > 0
	return
}

// 通过文件名称删除文件
func (h *Handler) RemoveObjectBatch(ctx context.Context, bucketName string, minioPaths []string) (err error) {
	for _, minioPath := range minioPaths {
		if err = h.RemoveObject(ctx, bucketName, minioPath); err != nil {
			return
		}
	}
	return
}
