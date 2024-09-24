package respx

import (
	"encoding/json"
	"net/url"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"github.com/tidwall/gjson"
)

// Header 导出表头
type Header struct {
	Key  string
	Name string
}

func ExcelHeaders(model any) (result []*Header) {
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

// BuildExcelByData 返回Excel二进制文件流
func BuildExcelByData(ctx *gin.Context, model any, data any, excelName string) {
	var xlsxFile = xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("Sheet1")
	if err != nil {
		Ctx(ctx).RespError(ExportFailedResp.SetData(err))
		ctx.Abort()
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
	excelName = url.QueryEscape(excelName)
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Add("Content-Disposition", "attachment;filename*=utf-8''"+excelName)
	ctx.Writer.Header().Add("Content-Transfer-Encoding", "binary")
	if err = xlsxFile.Write(ctx.Writer); err != nil {
		Ctx(ctx).RespError(ExportFailedResp.SetData(err))
		ctx.Abort()
		return
	}
	return
}
