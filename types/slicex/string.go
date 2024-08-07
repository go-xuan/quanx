package slicex

import (
	"sort"
	"strings"
)

// Split 字符串分割
func Split(str string, sep string) []string {
	slice := strings.Split(str, sep)
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}
	return slice
}

// Contains 字符串是否包含
func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// ContainsAny 数组是否包含
func ContainsAny(slice []string, args ...string) bool {
	if args != nil && len(args) > 0 {
		set := make(map[string]struct{})
		for _, item := range slice {
			set[item] = struct{}{}
		}
		for _, str := range args {
			if _, ok := set[str]; ok {
				return true
			}
		}
	}
	return false
}

// ContainsAll 数组是否包含
func ContainsAll(slice []string, args []string) (string, bool) {
	if args != nil && len(args) > 0 {
		set := make(map[string]struct{})
		for _, item := range slice {
			set[item] = struct{}{}
		}
		for _, str := range args {
			if _, ok := set[str]; !ok {
				return str, false
			}
		}
	}
	return "", true
}

// Distinct 合并并去重
func Distinct(slices ...[]string) (result []string) {
	set := make(map[string]struct{})
	for _, slice := range slices {
		for _, item := range slice {
			set[item] = struct{}{}
		}
	}
	for item := range set {
		result = append(result, item)
	}
	return result
}

// RetainAll 取交集
func RetainAll(slices ...[]string) (result []string) {
	set := make(map[string]int)
	for _, slice := range slices {
		for _, item := range slice {
			set[item]++
		}
	}
	for k, v := range set {
		if v > 1 {
			result = append(result, k)
		}
	}
	return result
}

// Exclude 移除
func Exclude(target []string, exclude []string) (result []string) {
	set := make(map[string]struct{})
	for _, item := range exclude {
		set[item] = struct{}{}
	}
	for _, item := range target {
		if _, ok := set[item]; !ok {
			result = append(result, item)
		}
	}
	return result
}

// SortAsc 数组正序
func SortAsc(slice []string) {
	sort.Sort(sort.StringSlice(slice))
}

// SortDesc 数组倒叙
func SortDesc(slice []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(slice)))
}
