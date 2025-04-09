package gormx

import (
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/utils/randx"
)

type Test struct {
	Id      string `json:"id" gorm:"type:string; comment:ID;"`
	Type    int    `json:"type" gorm:"type:int2; not null; comment:类型（1/2/3）"`
	Name    string `json:"name" gorm:"type:string; not null; comment:名字"`
	Address string `json:"address" gorm:"type:string; comment:地址"`
	FFF     string `json:"fff" orm:"type:string"`
}

func (t Test) TableName() string {
	return "quanx_test"
}

func (t Test) TableComment() string {
	return "测试"
}

func TestDatabase(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{
		Enable:   true,
		Type:     "mysql",
		Port:     3306,
		Username: "root",
		Password: "root",
		Database: "quanx",
	}, configx.FromDefault); err != nil {
		panic(err)
	}

	if err := InitTable("default", &Test{}); err != nil {
		fmt.Println(err)
	}

	GetInstance().Create(&Test{
		Id:   randx.UUID(),
		Type: randx.IntRange(1, 100),
		Name: randx.String(),
	})

	var tt2 = &Test{}
	GetInstance().First(tt2)
	fmt.Println(tt2)
}
