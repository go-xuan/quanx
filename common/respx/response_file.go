package respx

import (
	"encoding/json"
	"net/url"
	"path/filepath"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"github.com/tidwall/gjson"
)

// 导出表头
type Header struct {
	Key  string
	Name string
}

func ExcelHeaders(model interface{}) (result []*Header) {
	typeRef := reflect.TypeOf(model)
	for i := 0; i < typeRef.NumField(); i++ {
		if typeRef.Field(i).Tag.Get("export") != "" {
			result = append(result, &Header{
				Key:  typeRef.Field(i).Tag.Get("json"),
				Name: typeRef.Field(i).Tag.Get("export"),
			})
		}
	}
	return
}

// 返回Excel二进制文件流
func BuildExcelByData(ctx *gin.Context, model interface{}, data interface{}, fileName string) {
	var xlsxFile = xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("Sheet1")
	if err != nil {
		BuildException(ctx, ExportErr, err)
		return
	}
	// 写入表头
	headerRow := sheet.AddRow()
	headers := ExcelHeaders(model)
	for _, header := range headers {
		headerRow.AddCell().Value = header.Name
	}
	// 写入数据
	var dataBytes, _ = json.Marshal(data)
	var gjsonResults = gjson.ParseBytes(dataBytes).Array()
	for _, gjsonResult := range gjsonResults {
		dataMap := gjsonResult.Map()
		row := sheet.AddRow()
		for _, header := range headers {
			row.AddCell().Value = dataMap[header.Key].String()
		}
	}
	fileName = url.QueryEscape(fileName)
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Add("Content-Disposition", "attachment;filename*=utf-8''"+fileName)
	ctx.Writer.Header().Add("Content-Transfer-Encoding", "binary")
	err = xlsxFile.Write(ctx.Writer)
	if err != nil {
		BuildException(ctx, ExportErr, err)
		return
	}
	return
}

// Excel二进制文件流响应
func BuildExcelByFile(ctx *gin.Context, filePath string) {
	var xlsxFile, err = xlsx.OpenFile(filePath)
	if err != nil {
		BuildException(ctx, ExportErr, err)
		return
	}
	var fileName = url.QueryEscape(filepath.Base(filePath))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Add("Content-Disposition", "attachment;filename*=utf-8''"+fileName)
	ctx.Writer.Header().Add("Content-Transfer-Encoding", "binary")
	err = xlsxFile.Write(ctx.Writer)
	if err != nil {
		BuildException(ctx, ExportErr, err)
		return
	}
	return
}
