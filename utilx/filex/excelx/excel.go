package excelx

import (
	"path"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"

	"github.com/go-xuan/quanx/utilx/stringx"
)

type SheetInfoList []*SheetInfo
type SheetInfo struct {
	SheetName string `json:"sheetName"`
	StartRow  int    `json:"startRow"`
	EndRow    int    `json:"endRow"`
}

// 表头映射
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

func GetHeadersByStruct(st interface{}) (result []*Header) {
	typeRef := reflect.TypeOf(st)
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

// 在目标sheet页中，根据【具有合并单元格的行】进行横向拆分
// 取合【并单元格的值】作为拆分出的新sheet名
// excelPath:目标excel
// sheetName:目标sheet页,默认取第一个sheet页
func ExcelSplit(excelPath, sheetName string) (string, error) {
	var err error
	// 读取excel
	var xlsxFile *xlsx.File
	xlsxFile, err = xlsx.OpenFile(excelPath)
	if err != nil {
		return excelPath, err
	}
	// 读取目标sheet
	var theSheet *xlsx.Sheet
	if xlsxFile.Sheet[sheetName] == nil {
		sheetName = xlsxFile.Sheets[0].Name
	}
	for _, sheet := range xlsxFile.Sheets {
		if sheet.Name == sheetName {
			theSheet = sheet
		}
	}
	// 新增sheet页
	var NewExcelFile = xlsx.NewFile()
	var addSheetList SheetInfoList
	for rowNo, rowData := range theSheet.Rows {
		// 如果是合并单元格
		if rowData.Cells == nil || len(rowData.Cells) == 0 {
			continue
		} else if rowData.Cells[0].HMerge > 0 {
			addSheetName := rowData.Cells[0].Value
			if len(addSheetName) > 30 {
				addSheetName = addSheetName[:30]
			}
			_, err = NewExcelFile.AddSheet(addSheetName)
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
		if NewExcelFile.Sheet[item.SheetName] != nil {
			for rowNo, rowData := range theSheet.Rows {
				if rowNo >= item.StartRow && rowNo <= item.EndRow {
					row := NewExcelFile.Sheet[item.SheetName].AddRow()
					for _, cell := range rowData.Cells {
						row.AddCell().Value = cell.Value
					}
				}
			}
		}
	}
	dir, excelName := stringx.CutFromRight(excelPath, "\\")
	excelName, _ = stringx.CutFromLeft(excelName, ".")
	newExcelPath := path.Join(dir, excelName+"_split.xlsx")
	err = NewExcelFile.Save(newExcelPath)
	if err != nil {
		return excelPath, err
	}
	return newExcelPath, nil
}

// excelPath:目标excel
// sheetName:目标sheet页,默认取第一个sheet页
func ExcelReader(excelPath, sheetName string, headerMap map[string]string) ([]map[string]string, error) {
	var resultMapList []map[string]string
	var err error
	// 读取excel
	var xlsxFile *xlsx.File
	xlsxFile, err = xlsx.OpenFile(excelPath)
	if err != nil {
		return nil, err
	}
	// 读取目标sheet
	var theSheet *xlsx.Sheet
	if xlsxFile.Sheet[sheetName] == nil {
		theSheet = xlsxFile.Sheets[0]
	} else {
		theSheet = xlsxFile.Sheet[sheetName]
	}
	// 读取表头
	var headers []string
	for _, cell := range theSheet.Rows[0].Cells {
		header := headerMap[cell.Value]
		if header == "" {
			header = cell.Value
		}
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
					rowMap[headers[x]] = cell.String()
				} else {
					rowMap[headers[x]] = cell.String()
				}
			}
			resultMapList = append(resultMapList, rowMap)
		}
	}
	return resultMapList, nil
}

// 将数据写入excel
func ExcelWriter(excelPath string, st interface{}, dataList []map[string]string) (err error) {
	var xlsxFile = xlsx.NewFile()
	var sheet *xlsx.Sheet
	sheet, err = xlsxFile.AddSheet("Sheet1")
	if err != nil {
		log.Error("创建excel文件失败：%s", err)
		return
	}
	// 写入表头
	headerRow := sheet.AddRow()
	headers := GetHeadersByStruct(st)
	for _, header := range headers {
		headerRow.AddCell().Value = header.Name
	}
	// 写入数据
	for _, data := range dataList {
		row := sheet.AddRow()
		for _, header := range headers {
			row.AddCell().Value = data[header.Key]
		}
	}
	//这里从新生成
	err = xlsxFile.Save(excelPath)
	if err != nil {
		return
	}
	return
}
