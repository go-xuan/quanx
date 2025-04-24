package excelx

import (
	"strconv"
	"testing"
	"time"

	"github.com/tealeg/xlsx"

	"github.com/go-xuan/quanx/utils/randx"
)

func TestExcelWriter(t *testing.T) {
	var users []*User
	for i := 0; i < 10; i++ {
		users = append(users, &User{
			Name:     randx.Name(),
			Age:      strconv.Itoa(randx.IntRange(1, 100)),
			IdCard:   randx.IdCard(),
			Birthday: time.Now(),
		})
	}
	path := "test" + time.Now().Format("20060102150405") + ".xlsx"
	if err := WritToXlsx[*User](path, &User{}, users); err != nil {
		panic(err)
	}
}

type User struct {
	Name     string
	Age      string
	IdCard   string
	Birthday time.Time
}

func (u *User) AddHeader(sheet *xlsx.Sheet) {
	row := sheet.AddRow()
	row.AddCell().SetString("姓名")
	row.AddCell().SetString("年龄")
	row.AddCell().SetString("身份证号")
	row.AddCell().SetString("出生日期")
	row.AddCell().SetString("出生日期2")
}

func (u *User) AddRow(sheet *xlsx.Sheet) {
	row := sheet.AddRow()
	row.AddCell().SetString(u.Name)
	row.AddCell().SetString(u.Age)
	row.AddCell().SetString(u.IdCard)
	row.AddCell().SetDateTime(u.Birthday)
	row.AddCell().SetDateWithOptions(u.Birthday, xlsx.DateTimeOptions{
		Location:        time.Local,
		ExcelTimeFormat: "yyyy-mm-dd hh:mm:ss"},
	)
}
