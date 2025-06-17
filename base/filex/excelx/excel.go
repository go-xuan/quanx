package excelx

import (
	"encoding/json"
	"reflect"

	"github.com/tealeg/xlsx"
	"github.com/tidwall/gjson"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
)

const (
	TimeFmt = "yyyy-mm-dd h:mm:ss"
	DateFmt = "yyyy-mm-dd"
)

// GetHeaderMapping 通过反射获取excel表头映射，key:excel表头，value:结构体字段名
func GetHeaderMapping(v any) map[string]string {
	var result = make(map[string]string)
	var typeRef = reflect.TypeOf(v)
	for i := 0; i < typeRef.NumField(); i++ {
		field := typeRef.Field(i)
		if header := field.Tag.Get("excel"); header != "" {
			result[header] = field.Name
		}
	}
	return result
}

// GetSheet 获取sheet页
func GetSheet(path string, sheet string) (*xlsx.Sheet, error) {
	if file, err := xlsx.OpenFile(path); err != nil {
		return nil, errorx.Wrap(err, "open xlsx file error")
	} else if s, ok := file.Sheet[sheet]; ok {
		return s, nil
	} else {
		return file.Sheets[0], nil
	}
}

// ReadWithMapping 根据映射读取excel
func ReadWithMapping(path, sheet string, mapping map[string]string) ([]map[string]string, error) {
	// 读取目标sheet
	readSheet, err := GetSheet(path, sheet)
	if err != nil {
		return nil, errorx.Wrap(err, "get sheet error")
	}
	// 读取表头
	var headers []string
	for _, cell := range readSheet.Rows[0].Cells {
		header := cell.Value
		if mapping != nil && mapping[cell.Value] != "" {
			header = mapping[cell.Value]
		}
		headers = append(headers, header)
	}
	// 遍历excel(x:横向坐标，y:纵向坐标)
	var list = make([]map[string]string, 0)
	for y, row := range readSheet.Rows {
		if y > 0 {
			var data = make(map[string]string)
			for x, cell := range row.Cells {
				if x >= len(headers) {
					break
				}
				if cell.IsTime() {
					cell.SetFormat(TimeFmt)
				}
				data[headers[x]] = cell.String()
			}
			list = append(list, data)
		}
	}
	return list, nil
}

// ReadAny 根据结构体读取excel
func ReadAny[T any](path, sheet string, t T) ([]*T, error) {
	// 读取目标sheet
	readSheet, err := GetSheet(path, sheet)
	if err != nil {
		return nil, errorx.Wrap(err, "get sheet error")
	}
	if len(readSheet.Rows) == 0 {
		return nil, errorx.New("sheet is empty")
	}
	// 读取表头
	var headers []string
	var mapping = GetHeaderMapping(t)
	if len(mapping) == 0 {
		return nil, errorx.New("excel tag is required for struct")
	}
	for _, cell := range readSheet.Rows[0].Cells {
		headers = append(headers, stringx.IfZero(mapping[cell.Value], cell.Value))
	}
	// 遍历excel(x:横向坐标，y:纵向坐标)
	var list = make([]*T, 0)
	for y, row := range readSheet.Rows {
		if y > 0 {
			var data = make(map[string]string)
			for i, cell := range row.Cells {
				if i >= len(headers) {
					break
				}
				if cell.IsTime() {
					cell.SetFormat(TimeFmt)
				}
				data[headers[i]] = cell.String()
			}
			var item = new(T)
			if err = anyx.MapToStruct(data, item); err != nil {
				return nil, errorx.Wrap(err, "convert map to struct error")
			}
			list = append(list, item)
		}
	}
	return list, nil
}

// WriteAny 将数据写入excel
func WriteAny(path string, data any) error {
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
