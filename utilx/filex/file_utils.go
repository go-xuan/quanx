package filex

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-xuan/quanx/utilx/stringx"
)

const (
	DirAndFile = "all"
	OnlyDir    = "dir"
	OnlyFile   = "file"
	Json       = ".json"
	Yml        = ".yml"
	Yaml       = ".yaml"
	Toml       = ".toml"
	Properties = ".properties"
)

// 显示当前目录
func Pwd() string {
	_, fileStr, _, _ := runtime.Caller(1)
	dir, _ := filepath.Split(fileStr)
	return dir
}

// 拆分为文件路径和文件名
func SplitPath(path string) (dir string, file string) {
	if path == "" {
		return
	}
	if stringx.ContainsAny(path, "/", "\\") {
		dir, file = filepath.Split(path)
	} else {
		dir, file = "", path
	}
	return
}

// 判断是否文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径文件或文件夹是否存在
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// 创建文件
func Create(path string) {
	f, err := os.Create(path)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	if err != nil {
		log.Fatalf("db connect error: %#v\n", err.Error())
	}
}

// 创建文件夹
func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 先创建文件夹
		_ = os.MkdirAll(path, os.ModePerm)
		// 再修改权限
		_ = os.Chmod(path, os.ModePerm)
	}
}

// 获取后缀
func Suffix(path string) string {
	return filepath.Ext(filepath.Base(path))
}

// 获取文件名(不带后缀)
func FileName(path string) string {
	fullName := filepath.Base(path)
	return strings.TrimSuffix(fullName, filepath.Ext(fullName))
}

// 强制打开文件
func MustOpen(filePath string, fileName string) (*os.File, error) {
	fileAbsPath, err := filepath.Abs(filepath.Join(filePath, fileName))
	if err != nil {
		return nil, err
	}
	perm := CheckPermission(fileAbsPath)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", fileAbsPath)
	}
	err = NotExistCreateFile(fileAbsPath)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", fileAbsPath, err)
	}
	var file *os.File
	file, err = Open(fileAbsPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to OpenFile :%v", err)
	}
	return file, nil
}

// 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// 检查是否有权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// 不存在即创建
func NotExistCreateFile(src string) error {
	_, err := os.Stat(src)
	if notExist := os.IsNotExist(err); notExist == true {
		if _, err = os.Create(src); err != nil {
			return err
		}
	} else {
		if err = os.Remove(src); err == nil {
			if _, err = os.Create(src); err != nil {
				return err
			}
		}
	}
	return nil
}

// 通过url获取文件字节
func GetFileBytesByUrl(fileUrl string) ([]byte, error) {
	var result []byte
	var tr = &http.Transport{
		IdleConnTimeout:       time.Second * 2048,
		ResponseHeaderTimeout: time.Second * 10,
	}
	if strings.Index(fileUrl, "https") != -1 {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		requestURI, _ := url.ParseRequestURI(fileUrl)
		fileUrl = requestURI.String()
	}
	var client = &http.Client{Transport: tr}
	resp, err := client.Get(fileUrl)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)
	result, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 获取文件字节的二进制
func bytesToHexString(src []byte) string {
	if src == nil || len(src) == 0 {
		return ""
	}
	sb := strings.Builder{}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			sb.WriteString(strconv.FormatInt(int64(0), 10))
		}
		sb.WriteString(hv)
	}
	return sb.String()
}
