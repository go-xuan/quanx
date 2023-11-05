package filex

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-xuan/quanx/common/constx"
)

// 读取文件内容
func ReadFile(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if nil != err {
		return "", fmt.Errorf(" %s read file error: %v", filePath, err)
	}
	return string(bytes), nil
}

// 按行读取
func ReadFileLine(path string) (contents []string, err error) {
	var file *os.File
	file, err = os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
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
		contents = append(contents, string(line))
	}
	return contents, nil
}

// 更新文件
func ModifyFile(filePath string, replaces map[string]string) error {
	bytes, err := os.ReadFile(filePath)
	if nil != err {
		return fmt.Errorf(" %s read file error: %v", filePath, err)
	}
	content := string(bytes)
	for k, v := range replaces {
		content = strings.ReplaceAll(content, k, v)
	}
	return WriteFile(filePath, content, Overwrite)
}

type File struct {
	Path string
	Info os.FileInfo
}

// 获取目录下所有文件路径
func FileScan(dir string, ft string) (fileInfos []*File, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		file := File{path, info}
		switch ft {
		case constx.DirAndFile:
			fileInfos = append(fileInfos, &file)
		case constx.OnlyDir:
			if info.IsDir() {
				fileInfos = append(fileInfos, &file)
			}
		case constx.OnlyFile:
			if !info.IsDir() {
				fileInfos = append(fileInfos, &file)
			}
		default:
			if info.Name() == ft {
				fileInfos = append(fileInfos, &file)
			}
		}
		return nil
	})
	if err != nil {
		os.Exit(1)
	}
	return
}
