package crudx

import "gorm.io/gorm"

type Model[T any] struct {
	Struct T        // 表对应的结构体
	DB     *gorm.DB // gorm
}
