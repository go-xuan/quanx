package jsonx

import (
	"encoding/json"
	"os"
)

// 将json字符串转为interface{}
func JsonToInterface(jsonStr string, obj interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return err
	}
	return nil
}

// 结构体转Json
func InterfaceToJson(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 写入结构体到json文件
func writeStructToJson(path string, obj interface{}) {
	jsonBytes, err := json.MarshalIndent(obj, "", "	")
	if err != nil {
		return
	}
	err = os.WriteFile(path, jsonBytes, 0777)
	if err != nil {
		return
	}
}

// 读取json文件结构体
func readJsonToStruct(path string, obj interface{}) {
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonBytes, &obj)
	if err != nil {
		return
	}
}
