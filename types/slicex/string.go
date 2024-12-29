package slicex

import (
	"sort"
)

// Contains 字符串是否包含
func Contains[T string | int | int64 | float64](slice []T, v T) bool {
	for _, item := range slice {
		if item == v {
			return true
		}
	}
	return false
}

// ContainsAny 数组是否包含
func ContainsAny[T string | int | int64 | float64](slice []T, args ...T) bool {
	if args != nil && len(args) > 0 {
		m := make(map[T]struct{})
		for _, item := range slice {
			m[item] = struct{}{}
		}
		for _, arg := range args {
			if _, ok := m[arg]; ok {
				return true
			}
		}
	}
	return false
}

// ContainsAll 数组是否包含
func ContainsAll[T string | int | int64 | float64](slice []T, args ...T) bool {
	if args != nil && len(args) > 0 {
		var m = make(map[T]struct{})
		for _, k := range slice {
			m[k] = struct{}{}
		}
		for _, arg := range args {
			if _, ok := m[arg]; !ok {
				return false
			}
		}
	}
	return true
}

func Distinct[T string | int | int64 | float64](slices ...[]T) []T {
	if len(slices) > 0 {
		var m = make(map[T]struct{})
		for _, slice := range slices {
			for _, k := range slice {
				m[k] = struct{}{}
			}
		}
		var set []T
		for k := range m {
			set = append(set, k)
		}
		return set
	}
	return nil
}

// RetainAll 取交集
func RetainAll[T string | int | int64 | float64](slices ...[]T) []T {
	if len(slices) > 0 {
		var m = make(map[T]int)
		for _, slice := range slices {
			for _, k := range slice {
				m[k]++
			}
		}
		var result []T
		for k, v := range m {
			if v > 1 {
				result = append(result, k)
			}
		}
		return result
	}
	return nil
}

// Exclude 移除
func Exclude[T string | int | int64 | float64](target []T, exclude []T) []T {
	if len(target) > 0 && len(exclude) > 0 {
		m := make(map[T]struct{})
		for _, item := range exclude {
			m[item] = struct{}{}
		}
		var result []T
		for _, item := range target {
			if _, ok := m[item]; !ok {
				result = append(result, item)
			}
		}
		return result
	}
	return target
}

// SortAsc 数组正序
func SortAsc(slice []string) {
	sort.Sort(sort.StringSlice(slice))
}

// SortDesc 数组倒叙
func SortDesc(slice []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(slice)))
}
