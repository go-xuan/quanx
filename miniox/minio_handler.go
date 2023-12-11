package miniox

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/utilx/filex"
)

var instance *Handler

// minio控制器
type Handler struct {
	Config *Config       // minio配置
	Client *minio.Client // minio客户端
}

func This() *Handler {
	if instance == nil {
		panic("The minio instance has not been initialized, please check the relevant config")
	}
	return instance
}

// 创建桶
func (h *Handler) CreateBucket(ctx context.Context, bucketName string) error {
	isExist, err := h.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if isExist == false {
		err = h.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: ""})
		if err != nil {
			log.Warnf("创建桶失败 error , %+v", err)
			return err
		}
		err = h.Client.SetBucketPolicy(ctx, bucketName, defaultBucketPolicy(bucketName))
		if err != nil {
			log.Warnf("设置桶权限失败 error , %+v", err)
			return err
		}
	}
	return nil
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
	expires := time.Duration(h.Config.ExpireHour) * time.Hour
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
	fileBytes, err = filex.GetFileBytesByUrl(url)
	if err != nil {
		log.Warnf("读取文件失败 , %+v", err)
		return
	}
	minioPath = h.NewMinioPath(fileName)
	err = h.PutObject(ctx, bucketName, minioPath, bytes.NewBuffer(fileBytes))
	if err != nil {
		log.Warnf("上传文件失败 , %+v", err)
		return
	}
	return
}

// 上传文件
func (h *Handler) UploadFile(ctx context.Context, bucketName string, minioPath string, file *multipart.FileHeader) (err error) {
	var exist bool
	exist, err = h.ObjectExist(ctx, h.Config.BucketName, minioPath)
	if err != nil {
		log.Error("判断对象是否存在失败 error : ", err)
		return
	}
	if !exist {
		var mf multipart.File
		mf, err = file.Open()
		if err != nil {
			log.Warnf("打开文件失败 , %+v", err)
			return
		}
		defer func(mf multipart.File) {
			_ = mf.Close()
		}(mf)
		err = h.PutObject(ctx, bucketName, minioPath, mf)
		if err != nil {
			log.Warn("上传文件失败 , ", err)
			return
		}
	}
	return
}

// 获取对象是否存在
func (h *Handler) ObjectExist(ctx context.Context, bucketName string, minioPath string) (bool, error) {
	objInfo, err := h.Client.StatObject(ctx, bucketName, minioPath, minio.StatObjectOptions{})
	if err != nil {
		var minioError minio.ErrorResponse
		errByte, _ := json.Marshal(err)
		_ = json.Unmarshal(errByte, &minioError)
		if minioError.StatusCode == http.StatusNotFound || minioError.Code == "NoSuchKey" {
			return false, nil
		} else {
			log.Error("获取文件失败 , %+v", err)
			return false, err
		}
	}
	return objInfo.Size > 0, nil
}

// 通过文件名称删除文件
func (h *Handler) RemoveObjectBatch(ctx context.Context, bucketName string, minioPaths []string) error {
	for _, minioPath := range minioPaths {
		err := h.RemoveObject(ctx, bucketName, minioPath)
		if err != nil {
			log.Error("删除文件失败 , %+v", err)
			return err
		}
	}
	return nil
}
