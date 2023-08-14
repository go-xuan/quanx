package codex

import (
	"path/filepath"
	"strings"

	"github.com/quanxiaoxuan/quanx/utils/codex/template"
	"github.com/quanxiaoxuan/quanx/utils/filex"
	"github.com/quanxiaoxuan/quanx/utils/stringx"
)

const (
	DEMO       = "demo"
	TXT        = ".txt"
	Java       = ".java"
	Sql        = ".sql"
	GO         = ".go"
	Controller = "controller"
	Service    = "service"
	Dao        = "dao"
	Params     = "params"
	LineSep    = "\n"
	Tab        = "\t"
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
		name = DEMO
	}
	name = strings.ToLower(name)
	return name
}

// 生成通用代码（保留java类，go结构体，gorm结构体，以及增删改查sql）
func GenCodeByFieldList(saveDir, name string, fieldList FieldList) (outPath string, err error) {
	name = checkName(name)
	// 合并所有
	all := strings.Builder{}
	all.WriteString(LineSep)
	all.WriteString(BuildJavaClass(name, fieldList))
	all.WriteString(LineSep)
	all.WriteString(BuildGoStruct(name, fieldList))
	all.WriteString(LineSep)
	all.WriteString(BuildGormStruct(name, fieldList))
	all.WriteString(LineSep)
	all.WriteString(BuildSelectSql(name, fieldList))
	all.WriteString(LineSep)
	all.WriteString(BuildSelectSqlAlias(name, fieldList))
	all.WriteString(LineSep)
	all.WriteString(BuildInsertSql(name, fieldList))
	all.WriteString(LineSep)
	all.WriteString(BuildUpdateSql(name, fieldList))
	// 写入文件
	outPath = filepath.Join(saveDir, name+TXT)
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
	content.WriteString(LineSep)
	content.WriteString(BuildJavaClass(name, fieldList))
	content.WriteString(LineSep)
	// 写入文件
	outPath = filepath.Join(saveDir, stringx.UpperCamelCase(name)+Java)
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
	content.WriteString(LineSep)
	content.WriteString(BuildSelectSql(name, fieldList))
	content.WriteString(LineSep)
	content.WriteString(BuildSelectSqlAlias(name, fieldList))
	content.WriteString(LineSep)
	content.WriteString(BuildInsertSql(name, fieldList))
	content.WriteString(LineSep)
	content.WriteString(BuildUpdateSql(name, fieldList))
	// 写入文件
	outPath = filepath.Join(saveDir, name+Sql)
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
	content.WriteString(LineSep)
	content.WriteString(BuildGoStruct(name, fieldList))
	content.WriteString(LineSep)
	content.WriteString(LineSep)
	content.WriteString(BuildGormStruct(name, fieldList))
	content.WriteString(LineSep)
	// 写入文件
	outPath = filepath.Join(saveDir, name+GO)
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
	content.WriteString(LineSep)
	content.WriteString(BuildCkCreateSql(name, engine, fieldList))
	content.WriteString(LineSep)
	// 写入文件
	outPath = filepath.Join(saveDir, name+Sql)
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
	paramsFile := filepath.Join(saveDir, Params, name+"_"+Params+GO)
	controllerFile := filepath.Join(saveDir, Controller, name+"_"+Controller+GO)
	serviceFile := filepath.Join(saveDir, Service, name+"_"+Service+GO)
	daoFile := filepath.Join(saveDir, Dao, name+"_"+Dao+GO)
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
