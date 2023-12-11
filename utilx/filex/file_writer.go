package filex

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"os"

	"github.com/go-xuan/quanx/constx"
)

type writeMode uint

const (
	Overwrite writeMode = iota
	Append
)

// 写入文件
func WriteFile(path, content string, mode writeMode) error {
	dir, _ := SplitPath(path)
	CreateDir(dir)
	var flag int
	switch mode {
	case Overwrite:
		flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	case Append:
		flag = os.O_CREATE | os.O_RDWR | os.O_APPEND
	default:
		flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	file, err := os.OpenFile(path, flag, 0644)
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
func WriteFileLine(path string, contents []string, mode writeMode) error {
	var flag int
	switch mode {
	case Overwrite:
		flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	case Append:
		flag = os.O_CREATE | os.O_RDWR | os.O_APPEND
	default:
		flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	writer := bufio.NewWriter(file)
	for _, line := range contents {
		_, _ = writer.WriteString(line)
		_, _ = writer.WriteString(constx.NextLine)
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
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = ','
	csvWriter.UseCRLF = true
	if err = csvWriter.WriteAll(data); err != nil {
		return err
	}
	csvWriter.Flush()
	return nil
}
