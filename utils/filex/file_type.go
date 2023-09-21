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
	DirAndFile     = "all"
	OnlyDir        = "dir"
	OnlyFile       = "file"
	GoType         = ".go"
	ModType        = ".mod"
	JavaType       = ".java"
	JarType        = ".jar"
	ClassType      = ".class"
	SqlType        = ".sql"
	LogType        = ".log"
	ShellType      = ".sh"
	BatType        = ".bat"
	TxtType        = ".txt"
	JsonType       = ".json"
	YmlType        = ".yml"
	YamlType       = ".yaml"
	TomlType       = ".toml"
	XmlType        = ".xml"
	HtmlType       = ".html"
	PropertiesType = ".properties"
	DocType        = ".doc"
	DocxType       = ".docx"
	XlsType        = ".xls"
	XlsxType       = ".xlsx"
	PptType        = ".ppt"
	PptxType       = ".pptx"
	PdfType        = ".pdf"
	MdType         = ".md"
	Mp3Type        = ".mp3"
	Mp4Type        = ".mp4"
	JpgType        = ".jpg"
	HeicType       = ".heic"
	WavType        = ".wav"
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
