package gormx

import (
	"fmt"
	"testing"

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
	return "quanx_test"
}

func (t Test) InitData() any {
	return nil
}

func TestDatabase(t *testing.T) {
	// 先初始化redis
	if err := (&Config{
		Source:   "default",
		Enable:   true,
		Type:     "postgres",
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres",
		Database: "quanx",
		Debug:    true,
	}).Execute(); err != nil {
		fmt.Println(err)
	}
	if err := InitTable("default", &Test{}); err != nil {
		fmt.Println(err)
	}

	GetInstance().Model(Test{}).Create(&Test{
		Id:   randx.UUID(),
		Type: randx.IntRange(1, 100),
		Name: randx.Name(),
	})

	var tt2 = &Test{}
	GetInstance().Model(Test{}).First(tt2)

	fmt.Println(tt2)
}
