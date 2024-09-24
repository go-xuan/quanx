package gormx

import (
	"fmt"
	"testing"

	"github.com/go-xuan/quanx/core/configx"
)

type Test struct {
	Id   string `json:"id" gorm:"type:string; comment:ID;"`
	Type int    `json:"type" gorm:"type:int2; not null; comment:类型（1/2/3）"`
	Name string `json:"name" gorm:"type:string; not null; comment:名字"`
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
	if err := configx.Execute(&Config{
		Source:   "default",
		Enable:   true,
		Type:     "postgres",
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres",
		Database: "quanx",
		Debug:    true,
	},
	); err != nil {
		fmt.Println(err)
	}
	if err := InitTable("default", &Test{}); err != nil {
		fmt.Println(err)
	}
}
