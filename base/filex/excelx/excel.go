package excelx

import (
	"reflect"
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/tealeg/xlsx"
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

func GetSheet(path string, sheetName string) (*xlsx.Sheet, error) {
	if file, err := xlsx.OpenFile(path); err != nil {
		return nil, errorx.Wrap(err, "open xlsx file error")
	} else {
		return anyx.IfZero(file.Sheet[sheetName], file.Sheets[0]), nil
	}
}

// SplitXlsx 在目标sheet页中，根据【具有合并单元格的行】进行横向拆分，并取【合并单元格的值】作为拆分出的新sheet名
func SplitXlsx(path, sheetName string) (string, error) {
	// 读取excel
	readSheet, err := GetSheet(path, sheetName)
	if err != nil {
		return "", errorx.Wrap(err, "get sheetName error")
	}
	// 新增sheet页
	var newFile = xlsx.NewFile()
	var sheets []*Sheet
	for i, row := range readSheet.Rows {
		// 如果是合并单元格
		if len(row.Cells) > 0 && row.Cells[0].HMerge > 0 {
			name := row.Cells[0].Value
			if len(name) > 30 {
				name = name[:30]
			}
			sheet, _ := newFile.AddSheet(name)
			sheets = append(sheets, &Sheet{Name: name, StartRow: i, Sheet: sheet})
		} else {
			continue
		}
	}

	for i, sheet := range sheets {
		start, end := sheet.StartRow+1, len(readSheet.Rows)-1
		if i < len(sheets)-1 {
			end = sheets[i+1].StartRow - 2
		}
		for _, row := range readSheet.Rows[start:end] {
			newRow := sheet.Sheet.AddRow()
			for _, cell := range row.Cells {
				newRow.AddCell().Value = cell.Value
			}
		}
	}
	path = stringx.Insert(path, "_split", strings.LastIndex(path, ".")-1)
	if err = newFile.Save(path); err != nil {
		return "", errorx.Wrap(err, "save xlsx file error")
	}
	return path, nil
}

// ReadXlsxWithMapping 根据映射读取excel
func ReadXlsxWithMapping(path, sheetName string, mapping map[string]string) ([]map[string]string, error) {
	// 读取目标sheet
	readSheet, err := GetSheet(path, sheetName)
	if err != nil {
		return nil, errorx.Wrap(err, "get sheetName error")
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
	var data = make([]map[string]string, 0)
	for y, row := range readSheet.Rows {
		if y > 0 {
			var rowMap = make(map[string]string)
			for x, cell := range row.Cells {
				if x >= len(headers) {
					break
				}
				if cell.IsTime() {
					cell.SetFormat(timeFormat)
				}
				rowMap[headers[x]] = cell.String()
			}
			data = append(data, rowMap)
		}
	}
	return data, nil
}

// ReadXlsxWithStruct 根据结构体读取excel
func ReadXlsxWithStruct[T any](path, sheetName string, t T) ([]*T, error) {
	// 读取目标sheet
	readSheet, err := GetSheet(path, sheetName)
	if err != nil {
		return nil, errorx.Wrap(err, "get sheetName error")
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
	var data = make([]*T, 0)
	for y, row := range readSheet.Rows {
		if y > 0 {
			var rowMap = make(map[string]string)
			for i, cell := range row.Cells {
				if i >= len(headers) {
					break
				}
				if cell.IsTime() {
					cell.SetFormat("yyyy-mm-dd h:mm:ss")
				}
				rowMap[headers[i]] = cell.String()
			}
			item := new(T)
			if err = anyx.MapToStruct(rowMap, item); err != nil {
				return nil, errorx.Wrap(err, "convert map to struct error")
			}
			data = append(data, item)
		}
	}
	return data, nil
}
