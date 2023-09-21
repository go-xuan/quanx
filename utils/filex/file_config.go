package filex

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/BurntSushi/toml"
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
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		return
	}
	switch Suffix(filePath) {
	case JsonType:
		err = ReadJsonToPointer(bytes, config)
	case YamlType, YmlType:
		err = ReadYamlToPointer(bytes, config)
	case TomlType:
		err = ReadTomlToPointer(bytes, config)
	}
	return
}

// 读取json文件到指针
func ReadJsonToPointer(bytes []byte, config interface{}) (err error) {
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}

// 读取yaml文件到指针
func ReadYamlToPointer(bytes []byte, config interface{}) (err error) {
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}

// 读取toml文件到指针
func ReadTomlToPointer(bytes []byte, config interface{}) (err error) {
	err = toml.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	return
}
