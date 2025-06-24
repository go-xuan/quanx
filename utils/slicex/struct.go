package slicex

import "reflect"

// SliceToMapByField 切片映射（根据字段名）
func SliceToMapByField[T any](slice []T, field string) map[string]T {
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

// SliceToMapByIndex 切片映射（根据字段所在位置）
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

// SliceGroupByName 切片分组（根据字段名）
func SliceGroupByName[T any](slice []T, field string) map[string][]T {
	if slice == nil || len(slice) == 0 {
		return nil
	}
	if reflect.ValueOf(slice[0]).Kind() != reflect.Struct {
		panic("the type of the slice element must be struct")
	}
	var result = make(map[string][]T)
	for _, item := range slice {
		valueRef := reflect.ValueOf(item)
		key := valueRef.FieldByName(field).String()
		result[key] = append(result[key], item)
	}
	return result
}

// SliceGroupByIndex 切片分组（根据字段所在位置）
func SliceGroupByIndex[T any](slice []T, index int) map[string][]T {
	if slice == nil || len(slice) == 0 {
		return nil
	}
	if reflect.ValueOf(slice[0]).Kind() != reflect.Struct {
		panic("the type of the slice element must be struct")
	}
	var result = make(map[string][]T)
	for _, item := range slice {
		valueRef := reflect.ValueOf(item)
		key := valueRef.Field(index).String()
		result[key] = append(result[key], item)
	}
	return result
}

// PtrSliceGroupByField 指针切片分组（根据字段名）
func PtrSliceGroupByField[T any](PtrSlice []*T, field string) map[string][]*T {
	if PtrSlice == nil || len(PtrSlice) == 0 {
		return nil
	}
	if reflect.ValueOf(PtrSlice[0]).Kind() != reflect.Pointer {
		panic("the slice element kind must be pointer")
	}
	var result = make(map[string][]*T)
	for _, ptr := range PtrSlice {
		valueRef := reflect.ValueOf(*ptr)
		key := valueRef.FieldByName(field).String()
		result[key] = append(result[key], ptr)
	}
	return result
}

// PtrSliceGroupByIndex 指针切片分组（根据字段所在位置）
func PtrSliceGroupByIndex[T any](PtrSlice []*T, index int) map[string][]*T {
	if PtrSlice == nil || len(PtrSlice) == 0 {
		return nil
	}
	if reflect.ValueOf(PtrSlice[0]).Kind() != reflect.Pointer {
		panic("the slice element kind must be pointer")
	}
	var result = make(map[string][]*T)
	for _, ptr := range PtrSlice {
		valueRef := reflect.ValueOf(*ptr)
		key := valueRef.Field(index).String()
		result[key] = append(result[key], ptr)
	}
	return result
}
