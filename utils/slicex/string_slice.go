package slicex

import "sort"

// 字符串是否包含
func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// 数组是否包含
func ContainsAny(slice []string, args ...string) bool {
	if args != nil && len(args) > 0 {
		set := make(map[string]bool)
		for _, item := range slice {
			set[item] = true
		}
		for _, str := range args {
			if set[str] {
				return true
			}
		}
	}
	return false
}

// 数组是否包含
func ContainsAll(slice []string, args []string) (string, bool) {
	if args != nil && len(args) > 0 {
		set := make(map[string]bool)
		for _, item := range slice {
			set[item] = true
		}
		for _, str := range args {
			if !set[str] {
				return str, false
			}
		}
	}
	return "", true
}

// 数组去重
func Distinct(slice []string) []string {
	set := make(map[string]bool)
	for _, s := range slice {
		set[s] = true
	}
	distinct := make([]string, 0, len(set))
	for k := range set {
		distinct = append(distinct, k)
	}
	return distinct
}

// 合并
func UnionAll(slices ...[]string) []string {
	set := make(map[string]bool)
	for _, slice := range slices {
		for _, item := range slice {
			set[item] = true
		}
	}
	merge := make([]string, 0, len(set))
	for k := range set {
		merge = append(merge, k)
	}
	return merge
}

// 取交集
func RetainAll(slices ...[]string) []string {
	set := make(map[string]int)
	for _, slice := range slices {
		for _, item := range slice {
			set[item]++
		}
	}
	var result []string
	for k, v := range set {
		if v > 1 {
			result = append(result, k)
		}
	}
	return result
}

// 移除
func Remove(slice []string, strs ...string) []string {
	result := append([]string(nil), slice...)
	for _, str := range strs {
		var n int
		for _, v := range result {
			if v != str {
				result[n] = v
				n++
			}
		}
		result = result[:n]
	}
	return result
}

// 数组正序
func SortAsc(slice []string) {
	sort.Sort(sort.StringSlice(slice))
}

// 数组倒叙
func SortDesc(slice []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(slice)))
}
