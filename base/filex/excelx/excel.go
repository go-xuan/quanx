package excelx

import (
	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/tealeg/xlsx"
	"reflect"
)

const (
	timeFormat = "yyyy-mm-dd h:mm:ss"
	dateFmt    = "yyyy-mm-dd"
)

type Sheet struct {
	Name     string      `json:"name"`
	StartRow int         `json:"startRow"`
	EndRow   int         `json:"endRow"`
	Sheet    *xlsx.Sheet `json:"-"`
}

// GetHeadersByReflect 通过反射获取excel标签
func GetHeadersByReflect(v any) []*Header {
	var result []*Header
	var typeRef = reflect.TypeOf(v)
	for i := 0; i < typeRef.NumField(); i++ {
		field := typeRef.Field(i)
		if field.Tag.Get("excel") != "" {
			result = append(result, &Header{
				Key:  field.Tag.Get("json"),
				Name: field.Tag.Get("excel"),
			})
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

// ReadXlsxWithMapping 根据映射读取excel
func ReadXlsxWithMapping(path, sheet string, mapping map[string]string) ([]map[string]string, error) {
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
					cell.SetFormat(timeFormat)
				}
				data[headers[x]] = cell.String()
			}
			list = append(list, data)
		}
	}
	return list, nil
}

// ReadXlsxWithStruct 根据结构体读取excel
func ReadXlsxWithStruct[T any](path, sheet string, t T) ([]*T, error) {
	// 读取目标sheet
	readSheet, err := GetSheet(path, sheet)
	if err != nil {
		return nil, errorx.Wrap(err, "get sheet error")
	}
	// 读取表头
	var mapping = make(map[string]string)
	for _, header := range GetHeadersByReflect(t) {
		mapping[header.Name] = header.Key
	}
	var headers []string
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
					cell.SetFormat(timeFormat)
				}
				data[headers[i]] = cell.String()
			}
			item := new(T)
			if err = anyx.MapToStruct(data, item); err != nil {
				return nil, errorx.Wrap(err, "convert map to struct error")
			}
			list = append(list, item)
		}
	}
	return list, nil
}
