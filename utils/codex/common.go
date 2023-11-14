package codex

import (
	constx2 "github.com/go-xuan/quanx/public/constx"
	"github.com/go-xuan/quanx/utils/defaultx"
	"path/filepath"
	"strings"

	"github.com/go-xuan/quanx/utils/codex/template"
	"github.com/go-xuan/quanx/utils/filex"
	"github.com/go-xuan/quanx/utils/stringx"
)

// 字段配置
type FieldList []*Field
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
func GenCodeByFieldList(saveDir, name string, fieldList FieldList) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	all := strings.Builder{}
	all.WriteString(constx2.NextLine)
	all.WriteString(BuildJavaClass(name, fieldList))
	all.WriteString(constx2.NextLine)
	all.WriteString(BuildGoStruct(name, fieldList))
	all.WriteString(constx2.NextLine)
	all.WriteString(BuildGormStruct(name, fieldList))
	all.WriteString(constx2.NextLine)
	all.WriteString(BuildSelectSql(name, fieldList))
	all.WriteString(constx2.NextLine)
	all.WriteString(BuildSelectSqlAlias(name, fieldList))
	all.WriteString(constx2.NextLine)
	all.WriteString(BuildInsertSql(name, fieldList))
	all.WriteString(constx2.NextLine)
	all.WriteString(BuildUpdateSql(name, fieldList))
	// 写入文件
	outPath = filepath.Join(saveDir, name+filex.Txt)
	err = filex.WriteFile(outPath, all.String(), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成Java-Class
func GenJavaClassByFieldList(saveDir, name string, fieldList FieldList) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	content := strings.Builder{}
	content.WriteString(constx2.NextLine)
	content.WriteString(BuildJavaClass(name, fieldList))
	content.WriteString(constx2.NextLine)
	// 写入文件
	outPath = filepath.Join(saveDir, stringx.UpperCamelCase(name)+filex.Java)
	err = filex.WriteFile(outPath, content.String(), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成SQL样例
func GenSqlByFieldList(saveDir, name string, fieldList FieldList) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	content := strings.Builder{}
	content.WriteString(constx2.NextLine)
	content.WriteString(BuildSelectSql(name, fieldList))
	content.WriteString(constx2.NextLine)
	content.WriteString(BuildSelectSqlAlias(name, fieldList))
	content.WriteString(constx2.NextLine)
	content.WriteString(BuildInsertSql(name, fieldList))
	content.WriteString(constx2.NextLine)
	content.WriteString(BuildUpdateSql(name, fieldList))
	// 写入文件
	outPath = filepath.Join(saveDir, name+filex.Sql)
	err = filex.WriteFile(outPath, content.String(), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成go结构体
func GenGoStructByFieldList(saveDir, name string, fieldList FieldList) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	content := strings.Builder{}
	content.WriteString(constx2.NextLine)
	content.WriteString(BuildGoStruct(name, fieldList))
	content.WriteString(constx2.NextLine)
	content.WriteString(constx2.NextLine)
	content.WriteString(BuildGormStruct(name, fieldList))
	content.WriteString(constx2.NextLine)
	// 写入文件
	outPath = filepath.Join(saveDir, name+filex.Go)
	err = filex.WriteFile(outPath, content.String(), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成CK建表结构体
func GenCkCreateByFieldList(saveDir, name, engine string, fieldList FieldList) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	content := strings.Builder{}
	content.WriteString(constx2.NextLine)
	content.WriteString(BuildCkCreateSql(name, engine, fieldList))
	content.WriteString(constx2.NextLine)
	// 写入文件
	outPath = filepath.Join(saveDir, name+filex.Sql)
	err = filex.WriteFile(outPath, content.String(), filex.Overwrite)
	if err != nil {
		return
	}
	return
}

// 生成go模板
func GenGoTemplateByName(saveDir, name string) (err error) {
	name = checkName(name)
	modelName := stringx.UpperCamelCase(name)
	goFileName := name + filex.Go
	err = filex.WriteFile(
		filepath.Join(saveDir, "model", goFileName),
		strings.ReplaceAll(template.GoModel, `{modelName}`, modelName),
		filex.Overwrite)
	if err != nil {
		return
	}
	err = filex.WriteFile(
		filepath.Join(saveDir, "controller", goFileName),
		strings.ReplaceAll(template.GoController, `{modelName}`, modelName),
		filex.Overwrite)
	if err != nil {
		return err
	}
	err = filex.WriteFile(
		filepath.Join(saveDir, "logic", goFileName),
		strings.ReplaceAll(template.GoLogic, `{modelName}`, modelName),
		filex.Overwrite)
	if err != nil {
		return
	}
	err = filex.WriteFile(
		filepath.Join(saveDir, "dao", goFileName),
		strings.ReplaceAll(template.GoDao, `{modelName}`, modelName),
		filex.Overwrite)
	if err != nil {
		return
	}
	return
}
