package filex

import (
	"crypto/tls"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DirAndFile = "all"
	OnlyDir    = "dir"
	OnlyFile   = "file"
	Go         = ".go"
	Java       = ".java"
	Sql        = ".sql"
	Txt        = ".txt"
	Json       = ".json"
	Yml        = ".yml"
	Yaml       = ".yaml"
	Toml       = ".toml"
	Properties = ".properties"
	Jar        = ".jar"
	Mod        = ".mod"
	Class      = ".class"
	Log        = ".log"
	Shell      = ".sh"
	Bat        = ".bat"
	Xml        = ".xml"
	Html       = ".html"
	Doc        = ".doc"
	Docx       = ".docx"
	Xls        = ".xls"
	Xlsx       = ".xlsx"
	Ppt        = ".ppt"
	Pptx       = ".pptx"
	Pdf        = ".pdf"
	Md         = ".md"
	Mp3        = ".mp3"
	Mp4        = ".mp4"
	Jpg        = ".jpg"
	Png        = ".png"
	Heic       = ".heic"
	Wav        = ".wav"
)

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
		log.Error("获取图片失败：", err)
		return nil, err
	}
	body := resp.Body
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)
	result, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error("获取图片失败：", err)
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
