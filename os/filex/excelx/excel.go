package excelx

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/tealeg/xlsx"
	"github.com/tidwall/gjson"

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
		return nil, err
	} else {
		sheet := anyx.IfZero(xlsxFile.Sheet[sheetName], xlsxFile.Sheets[0])
		return sheet, nil
	}
}

// ExcelSplit 在目标sheet页中，根据【具有合并单元格的行】进行横向拆分
// 取合【并单元格的值】作为拆分出的新sheet名
// excelPath:目标excel
// sheetName:目标sheet页,默认取第一个sheet页
func ExcelSplit(excelPath, sheetName string) (newExcelPath string, err error) {
	// 读取excel
	var xlsxFile *xlsx.File
	if xlsxFile, err = xlsx.OpenFile(excelPath); err != nil {
		return
	}
	// 读取目标sheet
	theSheet := anyx.IfZero(xlsxFile.Sheet[sheetName], xlsxFile.Sheets[0])
	// 新增sheet页
	var newExcelFile = xlsx.NewFile()
	var addSheetList []*SheetInfo
	for rowNo, rowData := range theSheet.Rows {
		// 如果是合并单元格
		if rowData.Cells == nil || len(rowData.Cells) == 0 {
			continue
		} else if rowData.Cells[0].HMerge > 0 {
			addSheetName := rowData.Cells[0].Value
			if len(addSheetName) > 30 {
				addSheetName = addSheetName[:30]
			}
			_, err = newExcelFile.AddSheet(addSheetName)
			addSheetList = append(addSheetList, &SheetInfo{addSheetName, rowNo, 0})
		} else {
			continue
		}
	}
	for i, item := range addSheetList {
		item.StartRow = item.StartRow + 1
		if i < len(addSheetList)-1 {
			item.EndRow = addSheetList[i+1].StartRow - 2
		} else {
			item.EndRow = len(theSheet.Rows) - 1
		}
	}
	for _, item := range addSheetList {
		if newExcelFile.Sheet[item.SheetName] != nil {
			for rowNo, rowData := range theSheet.Rows {
				if rowNo >= item.StartRow && rowNo <= item.EndRow {
					row := newExcelFile.Sheet[item.SheetName].AddRow()
					for _, cell := range rowData.Cells {
						row.AddCell().Value = cell.Value
					}
				}
			}
		}
	}
	newExcelPath = stringx.Insert(excelPath, "_split", strings.LastIndex(excelPath, ".")-1)
	if err = newExcelFile.Save(newExcelPath); err != nil {
		return
	}
	return
}

// ExcelReader excelPath:目标excel
// sheetName:目标sheet页,默认取第一个sheet页
func ExcelReader(excelPath, sheetName string, mapping map[string]string) (data []map[string]string, err error) {
	// 读取目标sheet
	var theSheet *xlsx.Sheet
	if theSheet, err = GetSheet(excelPath, sheetName); err != nil {
		return
	}
	// 读取表头
	var headers []string
	for _, cell := range theSheet.Rows[0].Cells {
		var header = stringx.IfZero(mapping[cell.Value], cell.Value)
		headers = append(headers, header)
	}
	// 遍历excel(x:横向坐标，y:纵向坐标)
	data = make([]map[string]string, 0)
	for y, row := range theSheet.Rows {
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
	return
}

// ExcelReaderAny excelPath:目标excel
// sheetName:目标sheet页,默认取第一个sheet页
func ExcelReaderAny[T any](excelPath, sheetName string, obj T) (data []*T, err error) {
	// 读取目标sheet
	var theSheet *xlsx.Sheet
	if theSheet, err = GetSheet(excelPath, sheetName); err != nil {
		return
	}
	// 读取表头
	var mapping = ExcelTagMapping(obj)
	var headers []string
	for _, cell := range theSheet.Rows[0].Cells {
		var header = stringx.IfZero(mapping[cell.Value], cell.Value)
		headers = append(headers, header)
	}
	// 遍历excel(x:横向坐标，y:纵向坐标)
	for y, row := range theSheet.Rows {
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
				return
			}
			data = append(data, item)
		}
	}
	return
}

// ExcelWriter 将数据写入excel
func ExcelWriter(excelPath string, obj any, data any) (err error) {
	var xlsxFile = xlsx.NewFile()
	var sheet *xlsx.Sheet
	if sheet, err = xlsxFile.AddSheet("Sheet1"); err != nil {
		return
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
		return
	}
	return
}
