package codex

import (
	"path/filepath"
	"strings"

	"github.com/go-xuan/quanx/common/constx"
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
	if name == "" {
		name = "demo"
	}
	name = strings.ToLower(name)
	return name
}

// 生成通用代码（保留java类，go结构体，gorm结构体，以及增删改查sql）
func GenCodeByFieldList(saveDir, name string, fieldList FieldList) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	all := strings.Builder{}
	all.WriteString(constx.NextLine)
	all.WriteString(BuildJavaClass(name, fieldList))
	all.WriteString(constx.NextLine)
	all.WriteString(BuildGoStruct(name, fieldList))
	all.WriteString(constx.NextLine)
	all.WriteString(BuildGormStruct(name, fieldList))
	all.WriteString(constx.NextLine)
	all.WriteString(BuildSelectSql(name, fieldList))
	all.WriteString(constx.NextLine)
	all.WriteString(BuildSelectSqlAlias(name, fieldList))
	all.WriteString(constx.NextLine)
	all.WriteString(BuildInsertSql(name, fieldList))
	all.WriteString(constx.NextLine)
	all.WriteString(BuildUpdateSql(name, fieldList))
	// 写入文件
	outPath = filepath.Join(saveDir, name+constx.Txt)
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
	content.WriteString(constx.NextLine)
	content.WriteString(BuildJavaClass(name, fieldList))
	content.WriteString(constx.NextLine)
	// 写入文件
	outPath = filepath.Join(saveDir, stringx.UpperCamelCase(name)+constx.Java)
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
	content.WriteString(constx.NextLine)
	content.WriteString(BuildSelectSql(name, fieldList))
	content.WriteString(constx.NextLine)
	content.WriteString(BuildSelectSqlAlias(name, fieldList))
	content.WriteString(constx.NextLine)
	content.WriteString(BuildInsertSql(name, fieldList))
	content.WriteString(constx.NextLine)
	content.WriteString(BuildUpdateSql(name, fieldList))
	// 写入文件
	outPath = filepath.Join(saveDir, name+constx.Sql)
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
	content.WriteString(constx.NextLine)
	content.WriteString(BuildGoStruct(name, fieldList))
	content.WriteString(constx.NextLine)
	content.WriteString(constx.NextLine)
	content.WriteString(BuildGormStruct(name, fieldList))
	content.WriteString(constx.NextLine)
	// 写入文件
	outPath = filepath.Join(saveDir, name+constx.Go)
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
	content.WriteString(constx.NextLine)
	content.WriteString(BuildCkCreateSql(name, engine, fieldList))
	content.WriteString(constx.NextLine)
	// 写入文件
	outPath = filepath.Join(saveDir, name+constx.Sql)
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
	paramsFile := filepath.Join(saveDir, constx.Params, name+"_"+constx.Params+constx.Go)
	controllerFile := filepath.Join(saveDir, constx.Controller, name+"_"+constx.Controller+constx.Go)
	serviceFile := filepath.Join(saveDir, constx.Service, name+"_"+constx.Service+constx.Go)
	daoFile := filepath.Join(saveDir, constx.Dao, name+"_"+constx.Dao+constx.Go)
	if err = filex.WriteFile(paramsFile, strings.ReplaceAll(template.GoParam, `{modelName}`, modelName), filex.Overwrite); err != nil {
		return err
	}
	if err = filex.WriteFile(controllerFile, strings.ReplaceAll(template.GoController, `{modelName}`, modelName), filex.Overwrite); err != nil {
		return err
	}
	if err = filex.WriteFile(serviceFile, strings.ReplaceAll(template.GoService, `{modelName}`, modelName), filex.Overwrite); err != nil {
		return err
	}
	if err = filex.WriteFile(daoFile, strings.ReplaceAll(template.GoDao, `{modelName}`, modelName), filex.Overwrite); err != nil {
		return err
	}
	return
}
