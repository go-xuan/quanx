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

const (
	JsonTag  = "json"
	ExcelTag = "excel"
)

type SheetInfo struct {
	SheetName string `json:"sheetName"`
	StartRow  int    `json:"startRow"`
	EndRow    int    `json:"endRow"`
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
		if typeRef.Field(i).Tag.Get(ExcelTag) != "" {
			result = append(result, &Header{
				Key:  typeRef.Field(i).Tag.Get(JsonTag),
				Name: typeRef.Field(i).Tag.Get(ExcelTag),
			})
		}
	}
	return result
}

func ExcelTagMapping(v any) map[string]string {
	var result = make(map[string]string)
	var headers = GetHeadersByReflect(v)
	for _, header := range headers {
		result[header.Name] = header.Key
	}
	return result
}

func GetSheet(filePath string, sheetName string) (*xlsx.Sheet, error) {
	if xlsxFile, err := xlsx.OpenFile(filePath); err != nil {
		return nil, errorx.Wrap(err, "xlsx.OpenFile error")
	} else {
		sheet := anyx.IfZero(xlsxFile.Sheet[sheetName], xlsxFile.Sheets[0])
		return sheet, nil
	}
}

// ExcelSplit 在目标sheet页中，根据【具有合并单元格的行】进行横向拆分
// 取合【并单元格的值】作为拆分出的新sheet名
// excelPath:目标excel
// sheetName:目标sheet页,默认取第一个sheet页
func ExcelSplit(excelPath, sheetName string) (string, error) {
	// 读取excel
	xlsxFile, err := xlsx.OpenFile(excelPath)
	if err != nil {
		return "", errorx.Wrap(err, "xlsx.OpenFile error")
	}
	// 读取目标sheet
	sheet := anyx.IfZero(xlsxFile.Sheet[sheetName], xlsxFile.Sheets[0])
	// 新增sheet页
	var newFile = xlsx.NewFile()
	var addSheetList []*SheetInfo
	for rowNo, rowData := range sheet.Rows {
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
			item.EndRow = len(sheet.Rows) - 1
		}
	}
	for _, item := range addSheetList {
		if newFile.Sheet[item.SheetName] != nil {
			for rowNo, rowData := range sheet.Rows {
				if rowNo >= item.StartRow && rowNo <= item.EndRow {
					row := newFile.Sheet[item.SheetName].AddRow()
					for _, cell := range rowData.Cells {
						row.AddCell().Value = cell.Value
					}
				}
			}
		}
	}
	path := stringx.Insert(excelPath, "_split", strings.LastIndex(excelPath, ".")-1)
	if err = newFile.Save(path); err != nil {
		return "", errorx.Wrap(err, "xlsx.SaveFile error")
	}
	return path, nil
}

// ExcelReader excelPath:目标excel
// sheetName:目标sheet页,默认取第一个sheet页
func ExcelReader(excelPath, sheetName string, mapping map[string]string) ([]map[string]string, error) {
	// 读取目标sheet
	sheet, err := GetSheet(excelPath, sheetName)
	if err != nil {
		return nil, errorx.Wrap(err, "GetSheet error")
	}
	// 读取表头
	var headers []string
	for _, cell := range sheet.Rows[0].Cells {
		var header = stringx.IfZero(mapping[cell.Value], cell.Value)
		headers = append(headers, header)
	}
	// 遍历excel(x:横向坐标，y:纵向坐标)
	var data = make([]map[string]string, 0)
	for y, row := range sheet.Rows {
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

// ExcelReaderAny excelPath:目标excel
// sheetName:目标sheet页,默认取第一个sheet页
func ExcelReaderAny[T any](excelPath, sheetName string, obj T) ([]*T, error) {
	// 读取目标sheet
	sheet, err := GetSheet(excelPath, sheetName)
	if err != nil {
		return nil, errorx.Wrap(err, "GetSheet error")
	}
	// 读取表头
	var mapping = ExcelTagMapping(obj)
	var headers []string
	for _, cell := range sheet.Rows[0].Cells {
		var header = stringx.IfZero(mapping[cell.Value], cell.Value)
		headers = append(headers, header)
	}
	// 遍历excel(x:横向坐标，y:纵向坐标)
	var data = make([]*T, 0)
	for y, row := range sheet.Rows {
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
			var item *T
			if err = anyx.MapToStruct(rowMap, item); err != nil {
				return nil, errorx.Wrap(err, "anyx.MapToStruct error")
			}
			data = append(data, item)
		}
	}
	return data, nil
}

// ExcelWriter 将数据写入excel
func ExcelWriter(excelPath string, obj any, data any) error {
	var xlsxFile = xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("Sheet1")
	if err != nil {
		return errorx.Wrap(err, "xlsxFile.AddSheet error")
	}
	// 写入表头
	var headerRow = sheet.AddRow()
	var headers = GetHeadersByReflect(obj)
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
	if err = xlsxFile.Save(excelPath); err != nil {
		return errorx.Wrap(err, "xlsxFile.Save error")
	}
	return nil
}
