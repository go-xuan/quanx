package slicex

import (
	"reflect"
)

// 切片转map
func SliceToMapByName[T any](slice []T, field string) map[string]T {
	if slice == nil || len(slice) == 0 {
		return nil
	}
	var result = make(map[string]T)
	for _, item := range slice {
		valueRef := reflect.ValueOf(item)
		key := valueRef.FieldByName(field).String()
		result[key] = item
	}
	return result
}

// 切片转map
func SliceToMapByIndex[T any](slice []T, index int) map[string]T {
	if slice == nil || len(slice) == 0 {
		return nil
	}
	var result = make(map[string]T)
	for _, item := range slice {
		valueRef := reflect.ValueOf(item)
		key := valueRef.Field(index).String()
		result[key] = item
	}
	return result
}

// 切片根据字段名分组
func SliceGroupByName[T any](slice []T, field string) map[string][]T {
	if slice == nil || len(slice) == 0 {
		return nil
	}
	if reflect.ValueOf(slice[0]).Kind() != reflect.Struct {
		panic("切片元素类型必须为结构体(struct)")
	}
	var result = make(map[string][]T)
	for _, item := range slice {
		valueRef := reflect.ValueOf(item)
		key := valueRef.FieldByName(field).String()
		result[key] = append(result[key], item)
	}
	return result
}

// 切片根据字段下标分组
func SliceGroupByIndex[T any](slice []T, index int) map[string][]T {
	if slice == nil || len(slice) == 0 {
		return nil
	}
	if reflect.ValueOf(slice[0]).Kind() != reflect.Struct {
		panic("切片元素类型必须为结构体(struct)")
	}
	var result = make(map[string][]T)
	for _, item := range slice {
		valueRef := reflect.ValueOf(item)
		key := valueRef.Field(index).String()
		result[key] = append(result[key], item)
	}
	return result
}

// 指针切片根据字段名分组
func PtrSliceGroupByName[T any](PtrSlice []*T, field string) map[string][]*T {
	if PtrSlice == nil || len(PtrSlice) == 0 {
		return nil
	}
	if reflect.ValueOf(PtrSlice[0]).Kind() != reflect.Pointer {
		panic("切片元素类型必须为指针(Pointer)")
	}
	var result = make(map[string][]*T)
	for _, ptr := range PtrSlice {
		valueRef := reflect.ValueOf(*ptr)
		key := valueRef.FieldByName(field).String()
		result[key] = append(result[key], ptr)
	}
	return result
}

// 指针切片根据字段下标分组
func PtrSliceGroupByIndex[T any](PtrSlice []*T, index int) map[string][]*T {
	if PtrSlice == nil || len(PtrSlice) == 0 {
		return nil
	}
	if reflect.ValueOf(PtrSlice[0]).Kind() != reflect.Pointer {
		panic("切片元素类型必须为指针(Pointer)")
	}
	var result = make(map[string][]*T)
	for _, ptr := range PtrSlice {
		valueRef := reflect.ValueOf(*ptr)
		key := valueRef.Field(index).String()
		result[key] = append(result[key], ptr)
	}
	return result
}
