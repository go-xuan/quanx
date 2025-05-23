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
func WriteXlsx(path string, data any) error {
	// 写入数据
	var file = xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return errorx.Wrap(err, "add sheet error")
	}

	// 解析数据
	var bytes []byte
	if bytes, err = json.Marshal(data); err != nil {
		return errorx.Wrap(err, "data marshal error")
	}
	var list []gjson.Result
	if list = gjson.ParseBytes(bytes).Array(); len(list) == 0 {
		return errorx.Wrap(err, "data is empty")
	}

	// 获取表头
	var headers []string
	if first := list[0]; first.IsObject() {
		first.ForEach(func(key, value gjson.Result) bool {
			headers = append(headers, key.String())
			return true
		})
	}

	// 写入表头
	var headerRow = sheet.AddRow()
	for _, header := range headers {
		headerRow.AddCell().SetString(header)
	}

	// 写入数据
	for _, item := range list {
		row := sheet.AddRow()
		for _, header := range headers {
			row.AddCell().SetString(item.Get(header).String())
		}
	}
	// 保存xlsx
	if err = file.Save(path); err != nil {
		return errorx.Wrap(err, "save xlsx file error")
	}
	return nil
}
