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
	"github.com/quanxiaoxuan/quanx/utils/filex"
	log "github.com/sirupsen/logrus"
)

var CTL *Controller

// minio控制器
type Controller struct {
	Config *Config       // minio配置
	Client *minio.Client // minio客户端
}

// 初始化minio控制器
func Init(conf *Config) {
	client, err := conf.NewClient()
	if err == nil {
		CTL = &Controller{Config: conf, Client: client}
		log.Info("Minio连接成功!", conf.Format())
	} else {
		log.Error("Minio连接失败!", conf.Format())
		log.Error("error : ", err)
	}
}

// 创建桶
func (ctl *Controller) CreateBucket(ctx context.Context, bucketName string) error {
	isExist, err := ctl.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if isExist == false {
		err = ctl.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: ""})
		if err != nil {
			log.Warnf("创建桶失败 error , %+v", err)
			return err
		}
		err = ctl.Client.SetBucketPolicy(ctx, bucketName, defaultBucketPolicy(bucketName))
		if err != nil {
			log.Warnf("设置桶权限失败 error , %+v", err)
			return err
		}
	}
	return nil
}

// 生成minio存储路径
func (ctl *Controller) NewMinioPath(fileName string) string {
	return ctl.Config.MinioPath(fileName)
}

// 上传文件
func (ctl *Controller) PutObject(ctx context.Context, bucketName, minioPath string, reader io.Reader) (err error) {
	_, err = ctl.Client.PutObject(ctx, bucketName, minioPath, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return
	}
	return
}

// 下载文件
func (ctl *Controller) FGetObject(ctx context.Context, bucketName, minioPath, savePath string) (err error) {
	err = ctl.Client.FGetObject(ctx, bucketName, minioPath, savePath, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	return
}

// 删除文件
func (ctl *Controller) RemoveObject(ctx context.Context, bucketName, minioPath string) (err error) {
	err = ctl.Client.RemoveObject(ctx, bucketName, minioPath, minio.RemoveObjectOptions{GovernanceBypass: true})
	if err != nil {
		return
	}
	return
}

// 下载链接
func (ctl *Controller) PresignedGetObject(ctx context.Context, minioPath string) (minioUrl string, err error) {
	var URL *url.URL
	expires := time.Duration(ctl.Config.ExpireHour) * time.Hour
	URL, err = ctl.Client.PresignedGetObject(ctx, ctl.Config.BucketName, minioPath, expires, nil)
	if err != nil {
		return
	}
	minioUrl = URL.String()
	return
}

// 通过文件路径上传文件到桶
func (ctl *Controller) UploadFileByUrl(ctx context.Context, bucketName string, fileName string, url string) (minioPath string, err error) {
	var fileBytes []byte
	fileBytes, err = filex.GetFileBytesByUrl(url)
	if err != nil {
		log.Warnf("读取文件失败 , %+v", err)
		return
	}
	minioPath = ctl.NewMinioPath(fileName)
	err = ctl.PutObject(ctx, bucketName, minioPath, bytes.NewBuffer(fileBytes))
	if err != nil {
		log.Warnf("上传文件失败 , %+v", err)
		return
	}
	return
}

// 上传文件
func (ctl *Controller) UploadFile(ctx context.Context, bucketName string, minioPath string, file *multipart.FileHeader) (err error) {
	var exist bool
	exist, err = ctl.ObjectExist(ctx, ctl.Config.BucketName, minioPath)
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
		err = ctl.PutObject(ctx, bucketName, minioPath, mf)
		if err != nil {
			log.Warn("上传文件失败 , ", err)
			return
		}
	}
	return
}

// 获取对象是否存在
func (ctl *Controller) ObjectExist(ctx context.Context, bucketName string, minioPath string) (bool, error) {
	objInfo, err := ctl.Client.StatObject(ctx, bucketName, minioPath, minio.StatObjectOptions{})
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
func (ctl *Controller) RemoveObjectBatch(ctx context.Context, bucketName string, minioPaths []string) error {
	for _, minioPath := range minioPaths {
		err := ctl.RemoveObject(ctx, bucketName, minioPath)
		if err != nil {
			log.Error("删除文件失败 , %+v", err)
			return err
		}
	}
	return nil
}
