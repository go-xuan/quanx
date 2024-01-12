package filex

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
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
	file, err = os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	// 按行处理txt
	reader := bufio.NewReader(file)
	for {
		var line []byte
		line, _, err = reader.ReadLine()
		if err == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}
	return
}

// 更新文件
func ContentReplace(filePath string, replaces map[string]string) error {
	bytes, err := os.ReadFile(filePath)
	if nil != err {
		return fmt.Errorf(" %s read file error: %v", filePath, err)
	}
	content := string(bytes)
	for k, v := range replaces {
		content = strings.ReplaceAll(content, k, v)
	}
	return WriteFile(filePath, content)
}

// 写入文件
func WriteFile(filePath, content string, mode ...int) error {
	var flag = Overwrite
	if len(mode) > 0 {
		flag = mode[0]
	}
	file, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString(content)
	if err = writer.Flush(); err != nil {
		return err
	}
	return nil
}

// 数组按行写入文件
func WriteFileLine(filePath string, content []string, mode ...int) error {
	var flag = Overwrite
	if len(mode) > 0 {
		flag = mode[0]
	}
	file, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	writer := bufio.NewWriter(file)
	for _, line := range content {
		_, _ = writer.WriteString(line)
		_, _ = writer.WriteString("\n")
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	return nil
}

// 写入json文件
func WriteJson(path string, obj interface{}) error {
	jsonByte, err := json.MarshalIndent(obj, "", "	")
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, jsonByte, 0644); err != nil {
		return err
	}
	return nil
}

// 写入csv文件
func WriteCSV(path string, data [][]string) error {
	file, err := os.OpenFile(path, Overwrite, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	writer := csv.NewWriter(file)
	writer.Comma = ','
	writer.UseCRLF = true
	if err = writer.WriteAll(data); err != nil {
		return err
	}
	writer.Flush()
	return nil
}

type File struct {
	Path string
	Info os.FileInfo
}

// 获取目录下所有文件路径
func FileScan(dir string, suffix string) (files []*File, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
	})
	if err != nil {
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
