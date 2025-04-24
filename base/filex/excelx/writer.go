package excelx

import (
	"encoding/json"

	"github.com/tealeg/xlsx"
	"github.com/tidwall/gjson"

	"github.com/go-xuan/quanx/base/errorx"
)

type Writer interface {
	AddHeader(sheet *xlsx.Sheet)
	AddRow(sheet *xlsx.Sheet)
}

// WritToXlsx 将数据写入excel
func WritToXlsx[W Writer](path string, header W, rows []W) error {
	var file = xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return errorx.Wrap(err, "add sheet error")
	}
	// 写入表头和数据
	header.AddHeader(sheet)
	for _, row := range rows {
		row.AddRow(sheet)
	}
	// 保存xlsx
	if err = file.Save(path); err != nil {
		return errorx.Wrap(err, "save xlsx file error")
	}
	return nil
}

// Header 表头映射
type Header struct {
	Key  string
	Name string
}

// WriteXlsx 将数据写入excel
func WriteXlsx(path string, data any, headers []*Header) error {
	var xlsxFile = xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("Sheet1")
	if err != nil {
		return errorx.Wrap(err, "add sheet error")
	}
	// 写入表头
	var headerRow = sheet.AddRow()
	for _, header := range headers {
		headerRow.AddCell().SetString(header.Name)
	}
	// 写入数据
	var bytes, _ = json.Marshal(data)
	var array = gjson.ParseBytes(bytes).Array()
	for _, item := range array {
		var dataMap = item.Map()
		var row = sheet.AddRow()
		for _, header := range headers {
			row.AddCell().SetString(dataMap[header.Key].String())
		}
	}
	// 保存xlsx
	if err = xlsxFile.Save(path); err != nil {
		return errorx.Wrap(err, "save xlsx file error")
	}
	return nil
}
