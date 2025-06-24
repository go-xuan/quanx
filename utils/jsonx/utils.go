package jsonx

import (
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Get 通过path路径获取json对应部分
func Get(json string, path string) gjson.Result {
	result, keys := gjson.Parse(json), strings.Split(path, ".")
	for _, key := range keys {
		if index, ok := isArrayKey(key); ok {
			result = getArrayElement(result, index)
		} else {
			result = result.Get(key)
		}

		if !result.Exists() {
			return gjson.Result{}
		}
	}

	return result
}

// getArrayElement 根据索引获取数组元素，如果索引越界则返回空结果
func getArrayElement(result gjson.Result, index int) gjson.Result {
	if result.IsArray() {
		array := result.Array()
		if index >= 0 && index < len(array) {
			return array[index]
		}
	}
	return gjson.Result{}
}

// isArrayKey 判断键是否为数组索引格式，并返回索引值
// 例如：[1] 返回 1, true
// 例如：[a] 返回 -1, false
func isArrayKey(key string) (int, bool) {
	if n := len(key); n >= 3 && key[0] == '[' && key[n-1] == ']' {
		key = strings.TrimSpace(key[1 : n-1])
		if index, err := strconv.Atoi(key); err == nil {
			return index, true
		}
	}
	return -1, false
}
