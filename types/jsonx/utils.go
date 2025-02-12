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
		if index, ok := isArrayKey(key); ok && result.IsArray() {
			if array := result.Array(); len(array) > index {
				result = array[index]
			} else {
				result = array[0]
			}
			continue
		} else if child := result.Get(key); child.Exists() {
			result = child
			continue
		} else {
			return gjson.Result{}
		}
	}
	return result
}

func isArrayKey(key string) (int, bool) {
	if n := len(key); n >= 3 && key[0] == '[' && key[n-1] == ']' {
		key = strings.TrimSpace(key[1 : n-1])
		if index, err := strconv.Atoi(key); err == nil {
			return index, true
		}
	}
	return -1, false
}
