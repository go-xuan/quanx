package respx

import (
	"encoding/json"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"github.com/tidwall/gjson"
)

// 导出表头
type Header struct {
	Key  string
	Name string
}

// Gin框架返回Excel文件流
func BuildExcelExport(ctx *gin.Context, headers []*Header, dataList interface{}, fileName string) (err error) {
	var xlsxFile = xlsx.NewFile()
	var sheet *xlsx.Sheet
	sheet, err = xlsxFile.AddSheet("Sheet1")
	if err != nil {
		return
	}
	// 写入表头
	headerRow := sheet.AddRow()
	for _, header := range headers {
		headerRow.AddCell().Value = header.Name
	}
	// 写入数据
	dataBytes, _ := json.Marshal(dataList)
	gjsonResults := gjson.ParseBytes(dataBytes).Array()
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
	return xlsxFile.Write(ctx.Writer)
}
