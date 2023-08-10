package filex

import (
	"encoding/json"
	"os"

	"github.com/jinzhu/configor"
)

// 写入结构体到json文件
func WriteStructToJson(filePath string, obj interface{}) (err error) {
	var bytes []byte
	bytes, err = json.MarshalIndent(obj, "", "	")
	if err != nil {
		return
	}
	err = os.WriteFile(filePath, bytes, 0777)
	if err != nil {
		return
	}
	return
}

// 读取json文件到结构体
func ReadJsonToStruct(filePath string, obj interface{}) (err error) {
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, obj)
	if err != nil {
		return
	}
	return
}

// 读取配置文件到指针
func ReadConfigToPointer(filePath string, config interface{}) (err error) {
	switch Suffix(filePath) {
	case JsonType:
		err = ReadJsonToPointer(filePath, config)
	case YamlType, YmlType:
		err = ReadYamlToPointer(filePath, config)
	}
	return
}

// 读取json文件到指针
func ReadJsonToPointer(filePath string, config interface{}) (err error) {
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}

// 读取yaml文件到指针
func ReadYamlToPointer(filePath string, config interface{}) (err error) {
	if err = configor.New(&configor.Config{
		Debug:       true,
		Environment: configor.ENV(),
	}).Load(&config, filePath); err != nil {
		return
	}
	return
}
