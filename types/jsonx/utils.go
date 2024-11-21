package jsonx

import (
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Get 通过path路径获取json对应部分
// 和 gjson.Get() 方法不同
func Get(json string, path string) gjson.Result {
	result, keys := gjson.Parse(json), strings.Split(path, ".")
	for _, key := range keys {
		if child := result.Get(key); child.Exists() {
			result = child
			continue
		} else if index, ok := isArrayKey(key); ok && result.IsArray() {
			if array := result.Array(); len(array) > index {
				result = array[index]
			} else {
				result = array[0]
			}
			continue
		} else {
			result = gjson.Result{}
		}
	}
	return result
}

func isArrayKey(key string) (int, bool) {
	if n := len(key); n > 3 && key[0] == '[' && key[n-1] == ']' {
		if x, err := strconv.Atoi(key[1 : n-1]); err == nil {
			return x, true
		}
	}
	return 0, false
}
