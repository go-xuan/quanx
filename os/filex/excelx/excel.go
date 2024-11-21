package excelx

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/tealeg/xlsx"
	"github.com/tidwall/gjson"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
)

type SheetInfo struct {
	Name     string `json:"sheetName"`
	StartRow int    `json:"startRow"`
	EndRow   int    `json:"endRow"`
}

// Header 表头映射
type Header struct {
	Key  string
	Name string
}

func HeaderStyle() *xlsx.Style {
	return &xlsx.Style{
		Border:    xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"},
		Fill:      xlsx.Fill{FgColor: "FF92D050", PatternType: "solid"},
		Font:      xlsx.Font{Size: 14, Name: "微软雅黑", Charset: 134, Bold: true},
		Alignment: xlsx.Alignment{Horizontal: "center", Vertical: "center"},
	}
}

func DefaultStyle() *xlsx.Style {
	return &xlsx.Style{
		Border:    xlsx.Border{Left: "thin", Right: "thin", Top: "thin", Bottom: "thin"},
		Fill:      xlsx.Fill{PatternType: "none"},
		Font:      xlsx.Font{Size: 12, Name: "微软雅黑", Charset: 134},
		Alignment: xlsx.Alignment{Horizontal: "left", Vertical: "center"},
	}
}

// GetHeadersByReflect 通过反射获取excel标签
func GetHeadersByReflect(v any) []*Header {
	var result []*Header
	var typeRef = reflect.TypeOf(v)
	for i := 0; i < typeRef.NumField(); i++ {
		if typeRef.Field(i).Tag.Get("excel") != "" {
			result = append(result, &Header{
				Key:  typeRef.Field(i).Tag.Get("json"),
				Name: typeRef.Field(i).Tag.Get("excel"),
			})
		}
	}
	return result
}

func GetSheet(path string, sheet string) (*xlsx.Sheet, error) {
	if file, err := xlsx.OpenFile(path); err != nil {
		return nil, errorx.Wrap(err, "xlsx.OpenFile error")
	} else {
		return anyx.IfZero(file.Sheet[sheet], file.Sheets[0]), nil
	}
}

// SplitXlsx 在目标sheet页中，根据【具有合并单元格的行】进行横向拆分，并取【合并单元格的值】作为拆分出的新sheet名
func SplitXlsx(path, sheet string) (string, error) {
	// 读取excel
	xlsxFile, err := xlsx.OpenFile(path)
	if err != nil {
		return "", errorx.Wrap(err, "xlsx.OpenFile error")
	}
	// 读取目标sheet
	readSheet := anyx.IfZero(xlsxFile.Sheet[sheet], xlsxFile.Sheets[0])
	// 新增sheet页
	var newFile = xlsx.NewFile()
	var addSheetList []*SheetInfo
	for rowNo, rowData := range readSheet.Rows {
		// 如果是合并单元格
		if rowData.Cells == nil || len(rowData.Cells) == 0 {
			continue
		} else if rowData.Cells[0].HMerge > 0 {
			newSheetName := rowData.Cells[0].Value
			if len(newSheetName) > 30 {
				newSheetName = newSheetName[:30]
			}
			newFile.AddSheet(newSheetName)
			addSheetList = append(addSheetList, &SheetInfo{newSheetName, rowNo, 0})
		} else {
			continue
		}
	}
	for i, item := range addSheetList {
		item.StartRow = item.StartRow + 1
		if i < len(addSheetList)-1 {
			item.EndRow = addSheetList[i+1].StartRow - 2
		} else {
			item.EndRow = len(readSheet.Rows) - 1
		}
	}
	for _, item := range addSheetList {
		if newFile.Sheet[item.Name] != nil {
			for rowNo, rowData := range readSheet.Rows {
				if rowNo >= item.StartRow && rowNo <= item.EndRow {
					row := newFile.Sheet[item.Name].AddRow()
					for _, cell := range rowData.Cells {
						row.AddCell().Value = cell.Value
					}
				}
			}
		}
	}
	resultPath := stringx.Insert(path, "_split", strings.LastIndex(path, ".")-1)
	if err = newFile.Save(resultPath); err != nil {
		return "", errorx.Wrap(err, "xlsx.SaveFile error")
	}
	return resultPath, nil
}

// ReadXlsxWithMapping 根据映射读取excel
func ReadXlsxWithMapping(path, sheet string, mapping map[string]string) ([]map[string]string, error) {
	// 读取目标sheet
	readSheet, err := GetSheet(path, sheet)
	if err != nil {
		return nil, errorx.Wrap(err, "GetSheet error")
	}
	// 读取表头
	var headers []string
	for _, cell := range readSheet.Rows[0].Cells {
		var header = stringx.IfZero(mapping[cell.Value], cell.Value)
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
					cell.SetFormat("yyyy-mm-dd h:mm:ss")
				}
				rowMap[headers[x]] = cell.String()
			}
			data = append(data, rowMap)
		}
	}
	return data, nil
}

// ReadXlsxWithStruct 根据结构体读取excel
func ReadXlsxWithStruct[T any](path, sheet string, t T) ([]*T, error) {
	// 读取目标sheet
	readSheet, err := GetSheet(path, sheet)
	if err != nil {
		return nil, errorx.Wrap(err, "GetSheet error")
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
			var item *T
			if err = anyx.MapToStruct(rowMap, item); err != nil {
				return nil, errorx.Wrap(err, "anyx.MapToStruct error")
			}
			data = append(data, item)
		}
	}
	return data, nil
}

// WriteXlsx 将数据写入excel
func WriteXlsx(path string, data any, headers []*Header) error {
	var xlsxFile = xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("Sheet1")
	if err != nil {
		return errorx.Wrap(err, "xlsxFile.AddSheet error")
	}
	// 写入表头
	var headerRow = sheet.AddRow()
	for _, header := range headers {
		headerRow.AddCell().Value = header.Name
	}
	// 写入数据
	var bytes, _ = json.Marshal(data)
	var array = gjson.ParseBytes(bytes).Array()
	for _, item := range array {
		var dataMap = item.Map()
		var row = sheet.AddRow()
		for _, header := range headers {
			row.AddCell().Value = dataMap[header.Key].String()
		}
	}
	// 这里重新生成
	if err = xlsxFile.Save(path); err != nil {
		return errorx.Wrap(err, "xlsxFile.Save error")
	}
	return nil
}
