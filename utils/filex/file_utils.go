package filex

import (
	"fmt"
	"github.com/quanxiaoxuan/quanx/common/constx"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/quanxiaoxuan/quanx/utils/stringx"
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
	if stringx.ContainsAny(path, constx.ForwardSlash, constx.BackSlash) {
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

func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func NotExistCreateFile(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if _, err := os.Create(src); err != nil {
			return err
		}
	} else {
		if err := os.Remove(src); err == nil {
			if _, err = os.Create(src); err != nil {
				return err
			}
		}
	}
	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}
