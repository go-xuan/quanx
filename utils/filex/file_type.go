package filex

import (
	"bytes"
	"crypto/tls"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var fileTypeMap sync.Map

func init() {
	fileTypeMap.Store("ffd8ffe000104a464946", "image/jpg")          //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "image/png")          //PNG (png)
	fileTypeMap.Store("47494638396126026f01", "image/gif")          //GIF (gif)
	fileTypeMap.Store("49492a00227105008037", "image/tif")          //TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "image/bmp")          //16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "image/bmp")          //24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "image/bmp")          //256色位图(bmp)
	fileTypeMap.Store("41433130313500000000", "image/dwg")          //CAD (dwg)
	fileTypeMap.Store("3c21444f435459504520", "text/html")          //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c68746d6c3e0", "text/html")                 //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c21646f637479706520", "text/htm")           //HTM (htm)
	fileTypeMap.Store("48544d4c207b0d0a0942", "application/css")    //css
	fileTypeMap.Store("696b2e71623d696b2e71", "application/js")     //js
	fileTypeMap.Store("7b5c727466315c616e73", "text/rtf")           //Rich Text Format (rtf)
	fileTypeMap.Store("38425053000100000000", "application/psd")    //Photoshop (psd)
	fileTypeMap.Store("46726f6d3a203d3f6762", "application/eml")    //Email [Outlook Express 6] (eml)
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "application/msword") //MS Excel 注意：word、msi 和 excel的文件头一样
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "application/vsd")    //Visio 绘图
	fileTypeMap.Store("5374616E64617264204A", "application/mdb")    //MS Access (mdb)
	fileTypeMap.Store("252150532D41646F6265", "application/ps")
	fileTypeMap.Store("255044462d312e350d0a", "application/pdf")  //Adobe Acrobat (pdf)
	fileTypeMap.Store("2e524d46000000120001", "application/rmvb") //rmvb/rm相同
	fileTypeMap.Store("464c5601050000000900", "application/flv")  //flv与f4v相同
	fileTypeMap.Store("00000020667479706d70", "mp4")
	fileTypeMap.Store("49443303000000002176", "mp3")
	fileTypeMap.Store("000001ba210001000180", "mpg") //
	fileTypeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	fileTypeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	fileTypeMap.Store("52494646d07d60074156", "avi")
	fileTypeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	fileTypeMap.Store("504b0304140000000800", "zip")
	fileTypeMap.Store("526172211a0700cf9073", "rar")
	fileTypeMap.Store("235468697320636f6e66", "ini")
	fileTypeMap.Store("504b03040a0000000000", "jar")
	fileTypeMap.Store("4d5a9000030000000400", "exe")        //可执行文件
	fileTypeMap.Store("3c25402070616765206c", "jsp")        //jsp文件
	fileTypeMap.Store("4d616e69666573742d56", "mf")         //MF文件
	fileTypeMap.Store("3c3f786d6c2076657273", "xml")        //xml文件
	fileTypeMap.Store("494e5345525420494e54", "sql")        //xml文件
	fileTypeMap.Store("7061636b616765207765", "java")       //java文件
	fileTypeMap.Store("406563686f206f66660d", "bat")        //bat文件
	fileTypeMap.Store("1f8b0800000000000000", "gz")         //gz文件
	fileTypeMap.Store("6c6f67346a2e726f6f74", "properties") //bat文件
	fileTypeMap.Store("cafebabe0000002e0041", "class")      //bat文件
	fileTypeMap.Store("49545346030000006000", "chm")        //bat文件
	fileTypeMap.Store("04000000010000001300", "mxp")        //bat文件
	fileTypeMap.Store("504b0304140006000800", "docx")       //docx文件
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	fileTypeMap.Store("6431303a637265617465", "torrent")    //
	fileTypeMap.Store("6D6F6F76", "mov")                    //Quicktime (mov)
	fileTypeMap.Store("FF575043", "wpd")                    //WordPerfect (wpd)
	fileTypeMap.Store("CFAD12FEC5FD746F", "dbx")            //Outlook Express (dbx)
	fileTypeMap.Store("2142444E", "pst")                    //Outlook (pst)
	fileTypeMap.Store("AC9EBD8F", "qdf")                    //Quicken (qdf)
	fileTypeMap.Store("E3828596", "pwl")                    //Windows Password (pwl)
	fileTypeMap.Store("2E7261FD", "ram")                    //Real Audio (ram)
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

// 通过文件前面几个字节来判断文件类型
func GetFileType(fileBytes []byte) string {
	var fileType string
	fileCode := bytesToHexString(fileBytes)
	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

// 获取文件字节的二进制
func bytesToHexString(src []byte) string {
	if src == nil || len(src) == 0 {
		return ""
	}
	res := bytes.Buffer{}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}
