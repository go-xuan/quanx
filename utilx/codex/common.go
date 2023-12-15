package codex

import (
	"path/filepath"
	"strings"

	"github.com/go-xuan/quanx/utilx/defaultx"
	"github.com/go-xuan/quanx/utilx/filex"
	"github.com/go-xuan/quanx/utilx/stringx"
)

// 字段配置
type Field struct {
	Name    string `json:"name"`    // 字段名
	Origin  string `json:"origin"`  // 原始字段名
	Type    string `json:"type"`    // 字段类型
	Comment string `json:"comment"` // 备注
	Default string `json:"default"` // 默认值
}

// 表名检查
func checkName(name string) string {
	name = defaultx.String(name, "demo")
	name = strings.ToLower(name)
	return name
}

// 生成通用代码（保留java类，go结构体，gorm结构体，以及增删改查sql）
func GenCodeByFieldList(saveDir, name string, fieldList []*Field) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	all := strings.Builder{}
	all.WriteString(BuildJavaClass(name, fieldList))
	all.WriteString(BuildGoStruct(name, fieldList))
	all.WriteString(BuildGormStruct(name, fieldList))
	all.WriteString(BuildSelectSql(name, fieldList))
	all.WriteString(BuildSelectSqlAlias(name, fieldList))
	all.WriteString(BuildInsertSql(name, fieldList))
	all.WriteString(BuildUpdateSql(name, fieldList))
	// 写入文件
	outPath = filepath.Join(saveDir, name+".txt")
	err = filex.WriteFile(outPath, all.String(), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成Java-Class
func GenJavaClassByFieldList(saveDir, name string, fieldList []*Field) (outPath string, err error) {
	name = checkName(name)
	// 写入文件
	outPath = filepath.Join(saveDir, stringx.UpperCamelCase(name)+".java")
	err = filex.WriteFile(outPath, BuildJavaClass(name, fieldList), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成SQL样例
func GenSqlByFieldList(saveDir, name string, fieldList []*Field) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	content := strings.Builder{}
	content.WriteString(BuildSelectSql(name, fieldList))
	content.WriteString(BuildSelectSqlAlias(name, fieldList))
	content.WriteString(BuildInsertSql(name, fieldList))
	content.WriteString(BuildUpdateSql(name, fieldList))
	// 写入文件
	outPath = filepath.Join(saveDir, name+".sql")
	err = filex.WriteFile(outPath, content.String(), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成go结构体
func GenGoStructByFieldList(saveDir, name string, fieldList []*Field) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	content := strings.Builder{}
	content.WriteString(BuildGoStruct(name, fieldList))
	content.WriteString(BuildGormStruct(name, fieldList))
	// 写入文件
	outPath = filepath.Join(saveDir, name+".go")
	err = filex.WriteFile(outPath, content.String(), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成CK建表结构体
func GenCkCreateByFieldList(saveDir, name, engine string, fieldList []*Field) (outPath string, err error) {
	name = checkName(name)
	// 写入文件
	outPath = filepath.Join(saveDir, name+".sql")
	err = filex.WriteFile(outPath, BuildCkCreateSql(name, engine, fieldList), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成go模板
func GenGoTemplateByName(saveDir, name string) (err error) {
	name = checkName(name)
	modelName := stringx.UpperCamelCase(name)
	fileName := name + ".go"
	err = filex.WriteFile(
		filepath.Join(saveDir, "model", fileName),
		strings.ReplaceAll(StructTemplate, `{modelName}`, modelName),
		filex.Overwrite)
	if err != nil {
		return
	}
	err = filex.WriteFile(
		filepath.Join(saveDir, "controller", fileName),
		strings.ReplaceAll(ControllerTemplate, `{modelName}`, modelName),
		filex.Overwrite)
	if err != nil {
		return err
	}
	err = filex.WriteFile(
		filepath.Join(saveDir, "logic", fileName),
		strings.ReplaceAll(LogicTemplate, `{modelName}`, modelName),
		filex.Overwrite)
	if err != nil {
		return
	}
	err = filex.WriteFile(
		filepath.Join(saveDir, "dao", fileName),
		strings.ReplaceAll(DaoTemplate, `{modelName}`, modelName),
		filex.Overwrite)
	if err != nil {
		return
	}
	return
}
