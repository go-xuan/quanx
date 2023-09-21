package crudx

import (
	log "github.com/sirupsen/logrus"
)

// 表结构
func (m *Model[T]) DbInitTable() (err error) {
	if m.DB.Migrator().HasTable(&m.Struct) {
		return m.DB.Migrator().AutoMigrate(&m.Struct)
	} else {
		return m.DB.Migrator().CreateTable(&m.Struct)
	}
}

// 新增
func (m *Model[T]) DbAdd(add T) (err error) {
	err = m.DB.Create(&add).Error
	if err != nil {
		log.Error("对象新增失败 ： ", err)
		return
	}
	return
}

// 更新
func (m *Model[T]) DbUpdate(update T) (err error) {
	err = m.DB.Updates(&update).Error
	if err != nil {
		log.Error("对象更新失败 ： ", err)
		return
	}
	return
}

// 删除
func (m *Model[T]) DbDelete(id string) (err error) {
	err = m.DB.Delete(&m.Struct, id).Error
	if err != nil {
		log.Error("对象删除失败 ： ", err)
		return
	}
	return
}

// 查询
func (m *Model[T]) DbDetail(id string) (detail T, err error) {
	err = m.DB.Find(&m.Struct, id).Scan(&detail).Error
	if err != nil {
		log.Error("对象查询失败 ： ", err)
		return
	}
	return
}
