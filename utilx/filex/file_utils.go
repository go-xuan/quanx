package filex

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-xuan/quanx/utilx/anyx"
	"github.com/go-xuan/quanx/utilx/stringx"
)

const (
	WriteOnly = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	Overwrite = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	Append    = os.O_RDWR | os.O_CREATE | os.O_APPEND
)

const (
	DirAndFile = "all"
	OnlyDir    = "dir"
	OnlyFile   = "file"
	Json       = "json"
	Yml        = "yml"
	Yaml       = "yaml"
	Toml       = "toml"
	Properties = "properties"
)

// 读取文件内容
func ReadFile(filePath string) (bytes []byte, err error) {
	bytes, err = os.ReadFile(filePath)
	if nil != err {
		err = fmt.Errorf(" %s read file error: %v", filePath, err)
		return
	}
	return
}

// 按行读取
func ReadFileLine(filePath string) (lines []string, err error) {
	var file *os.File
	if file, err = os.OpenFile(filePath, os.O_RDONLY, 0666); err != nil {
		return
	}
	defer file.Close()
	// 按行处理txt
	reader := bufio.NewReader(file)
	for {
		var line []byte
		if line, _, err = reader.ReadLine(); err == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}
	return
}

// 更新文件
func ContentReplace(filePath string, replaces map[string]string) (err error) {
	var bytes []byte
	if bytes, err = os.ReadFile(filePath); err != nil {
		return
	}
	content := string(bytes)
	for k, v := range replaces {
		content = strings.ReplaceAll(content, k, v)
	}
	return WriteFile(filePath, content)
}

// 写入文件
func WriteFile(filePath, content string, mode ...int) (err error) {
	CreateDirNotExist(filePath)
	var flag = anyx.If(len(mode) > 0, mode[0], Overwrite)
	var file *os.File
	if file, err = os.OpenFile(filePath, flag, 0644); err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	if _, err = writer.WriteString(content); err != nil {
		return
	}
	if err = writer.Flush(); err != nil {
		return
	}
	return
}

// 数组按行写入文件
func WriteFileLine(filePath string, content []string, mode ...int) (err error) {
	var flag = anyx.If(len(mode) > 0, mode[0], Overwrite)
	var file *os.File
	if file, err = os.OpenFile(filePath, flag, 0644); err != nil {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, line := range content {
		_, _ = writer.WriteString(line)
		_, _ = writer.WriteString("\n")
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	return
}

// 写入json文件
func WriteJson(filePath string, obj interface{}) (err error) {
	var bytes []byte
	if bytes, err = json.MarshalIndent(obj, "", "	"); err != nil {
		return
	}
	CreateDirNotExist(filePath)
	var file *os.File
	if file, err = os.OpenFile(filePath, Overwrite, 0644); err != nil {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	if _, err = writer.Write(bytes); err != nil {
		return
	}
	if err = writer.Flush(); err != nil {
		return
	}
	return
}

// 写入csv文件
func WriteCSV(filePath string, data [][]string) (err error) {
	CreateDirNotExist(filePath)
	var file *os.File
	if file, err = os.OpenFile(filePath, Overwrite, 0644); err != nil {
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ','
	writer.UseCRLF = true
	if err = writer.WriteAll(data); err != nil {
		return
	}
	writer.Flush()
	return
}

type File struct {
	Path string
	Info os.FileInfo
}

// 获取目录下所有文件路径
func FileScan(dir string, suffix string) (files []*File, err error) {
	if err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		file := File{path, info}
		switch suffix {
		case DirAndFile:
			files = append(files, &file)
		case OnlyDir:
			if info.IsDir() {
				files = append(files, &file)
			}
		case OnlyFile:
			if !info.IsDir() {
				files = append(files, &file)
			}
		default:
			if info.Name() == suffix {
				files = append(files, &file)
			}
		}
		return nil
	}); err != nil {
		os.Exit(1)
	}
	return
}

// 显示当前目录
func Pwd() string {
	_, fileStr, _, _ := runtime.Caller(1)
	dir, _ := filepath.Split(fileStr)
	return dir
}

// 拆分为文件路径和文件名
func SplitPath(path string) (dir string, fileName string) {
	if path == "" {
		return
	}
	if stringx.ContainsAny(path, "/", "\\") {
		dir, fileName = filepath.Split(path)
	} else {
		dir, fileName = "", path
	}
	return
}

// 判断是否文件夹
func IsDir(path string) bool {
	if s, err := os.Stat(path); err != nil {
		return false
	} else {
		return s.IsDir()
	}
}

// 判断所给路径文件或文件夹是否存在
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// 创建文件
func Create(path string) (err error) {
	var f *os.File
	if f, err = os.Create(path); err != nil {
		return
	}
	defer f.Close()
	return
}

// 创建文件
func CreateDirNotExist(filePath string) {
	dir, _ := filepath.Split(filePath)
	if !Exists(dir) {
		CreateDir(dir)
	}
}

// 创建文件夹
func CreateDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 先创建文件夹
		_ = os.MkdirAll(dir, os.ModePerm)
		// 再修改权限
		_ = os.Chmod(dir, os.ModePerm)
	}
}

// 获取后缀
func Suffix(path string) string {
	if path != "" {
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '.' {
				return path[i+1:]
			}
		}
	}
	return ""
}

// 获取文件名(不带后缀)
func FileName(path string) string {
	var fullName = filepath.Base(path)
	return strings.TrimSuffix(fullName, filepath.Ext(fullName))
}

// 强制打开文件
func MustOpen(filePath string, fileName string) (file *os.File, err error) {
	var fileAbsPath string
	if fileAbsPath, err = filepath.Abs(filepath.Join(filePath, fileName)); err != nil {
		return
	}
	if perm := CheckPermission(fileAbsPath); perm == true {
		err = fmt.Errorf("file permission denied : %s", fileAbsPath)
		return
	}
	if err = CreateIsNotExist(fileAbsPath); err != nil {
		return
	}

	if file, err = Open(fileAbsPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return
	}
	return
}

// 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	if f, err := os.OpenFile(name, flag, perm); err != nil {
		return nil, err
	} else {
		return f, nil
	}
}

// 检查是否有权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// 不存在即创建
func CreateIsNotExist(src string) error {
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
func GetFileBytesByUrl(fileUrl string) (bytes []byte, err error) {
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
	var resp *http.Response
	if resp, err = client.Get(fileUrl); err != nil {
		return
	}
	var body = resp.Body
	defer body.Close()
	if bytes, err = io.ReadAll(body); err != nil {
		return
	}
	return
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
