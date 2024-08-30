package filex

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/types/stringx"
)

const (
	Overwrite  = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	Append     = os.O_RDWR | os.O_CREATE | os.O_APPEND
	WriteOnly  = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	AppendOnly = os.O_WRONLY | os.O_CREATE | os.O_APPEND
)

const (
	DirAndFile = "all"
	OnlyDir    = "dir"
	OnlyFile   = "file"
)

// ReadFile 读取文件内容
func ReadFile(filePath string) (bytes []byte, err error) {
	if bytes, err = os.ReadFile(filePath); err != nil {
		return
	}
	return
}

// ReadFileLine 按行读取
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

// ContentReplace 更新文件
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

// WriteFile 写入文件
func WriteFile(filePath, content string, mode ...int) (err error) {
	CreateIfNotExist(filePath)
	var flag = intx.Default(Overwrite, mode...)
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

// FileSplit 文件拆分
func FileSplit(filePath string, size int) (paths []string, err error) {
	var file *os.File
	defer file.Close()
	if file, err = os.OpenFile(filePath, os.O_RDONLY, 0666); err != nil {
		return
	}
	dir, filename, suffix := Analyse(filePath)
	dir = filepath.Join(dir, filename)
	reader := bufio.NewReader(file)
	count, index := 1, 1
	var bf = strings.Builder{}
	for {
		if index < size {
			var line []byte
			if line, _, err = reader.ReadLine(); err == io.EOF {
				path := filepath.Join(dir, "split_"+strconv.Itoa(count)+suffix)
				if err = WriteFile(path, bf.String()); err != nil {
					return
				}
				paths = append(paths, path)
				break
			}
			bf.WriteString("\n")
			bf.Write(line)
		} else {
			index = 1
			path := filepath.Join(dir, "split_"+strconv.Itoa(count)+suffix)
			if err = WriteFile(path, bf.String()); err != nil {
				return
			}
			paths = append(paths, path)
			bf.Reset()
			count++
		}
		index++
	}
	return
}

// WriteFileLine 数组按行写入文件
func WriteFileLine(filePath string, content []string, mode ...int) (err error) {
	var flag = intx.Default(Overwrite, mode...)
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

// WriteJson 写入json文件
func WriteJson(filePath string, v any) (err error) {
	var bytes []byte
	if bytes, err = json.MarshalIndent(v, "", "	"); err != nil {
		return
	}
	CreateIfNotExist(filePath)
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

// WriteCSV 写入csv文件
func WriteCSV(filePath string, data [][]string) (err error) {
	CreateIfNotExist(filePath)
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

// FileScan 获取目录下所有文件路径
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

// Pwd 获取绝对路径
func Pwd(path ...string) (pwd string) {
	if len(path) == 0 {
		_, file, _, _ := runtime.Caller(1)
		pwd, _ = filepath.Split(file)
	} else {
		pwd, _ = filepath.Abs(path[0])
	}
	return
}

// SplitPath 拆分为文件路径和文件名
func SplitPath(path string) (dir, file string) {
	if path != "" {
		if stringx.ContainsAny(path, "/", "\\") {
			dir, file = filepath.Split(path)
		} else {
			file = path
		}
	}
	return
}

// IsDir 判断是否文件夹
func IsDir(path string) bool {
	if fileInfo, err := os.Stat(path); err != nil {
		return false
	} else {
		return fileInfo.IsDir()
	}
}

// Exists 判断所给路径文件或文件夹是否存在
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// Create 创建文件
func Create(path string) (err error) {
	var f *os.File
	if f, err = os.Create(path); err != nil {
		return
	}
	defer f.Close()
	return
}

// CreateIfNotExist 创建文件
func CreateIfNotExist(path string) {
	dir, _ := filepath.Split(path)
	CreateDir(dir)
	_ = Create(path)
}

// CreateDir 创建文件夹
func CreateDir(path string) {
	if !Exists(path) {
		dir, file := filepath.Split(path)
		if stringx.Index(file, ".") == -1 {
			dir = filepath.Join(dir, file)
		}
		// 先创建文件夹
		_ = os.MkdirAll(dir, os.ModePerm)
		// 再修改权限
		_ = os.Chmod(dir, os.ModePerm)
	}
}

// IsEmptyDir 检查给定的目录是否为空
func IsEmptyDir(dir string) bool {
	f, err := os.Open(dir)
	defer f.Close()
	if err != nil {
		return false
	}
	// 读取目录内容
	var names []string
	if names, err = f.Readdirnames(0); err != nil {
		return false
	}
	// 如果目录内容为空，则目录为空
	return len(names) == 0
}

// Analyse 获取文件夹、文件名，文件后缀
func Analyse(path string) (dir, name, suffix string) {
	if dir, name = filepath.Split(path); name != "" {
		for i := len(name) - 1; i >= 0; i-- {
			if name[i] == '.' {
				name, suffix = name[:i], name[i:]
			}
		}
	}
	return
}

// Suffix 获取后缀
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

// FileName 获取文件名(不带后缀)
func FileName(path string) string {
	var fullName = filepath.Base(path)
	return strings.TrimSuffix(fullName, filepath.Ext(fullName))
}

// MustOpen 强制打开文件
func MustOpen(filePath string, fileName string) (file *os.File, err error) {
	var fileAbsPath string
	if fileAbsPath, err = filepath.Abs(filepath.Join(filePath, fileName)); err != nil {
		return
	}
	if perm := CheckPermission(fileAbsPath); perm == true {
		err = errorx.Errorf("file permission denied: %s", fileAbsPath)
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

// Open 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	if f, err := os.OpenFile(name, flag, perm); err != nil {
		return nil, err
	} else {
		return f, nil
	}
}

// CheckPermission 检查是否有权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// CreateIsNotExist 不存在即创建
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

// GetFileBytesByUrl 通过url获取文件字节
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
