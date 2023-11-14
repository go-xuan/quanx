package crudx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/public/modelx"
	respx2 "github.com/go-xuan/quanx/public/respx"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/public/gormx"
)

func AddCrudRouter[T any](router *gin.RouterGroup, relativePath string, source ...string) {
	var crud = &Model[T]{DB: gormx.CTL.GetDB(source...)}
	group := router.Group(relativePath)
	group.POST("create", crud.Create) // 新增
	group.POST("update", crud.Update) // 修改
	group.GET("delete", crud.Delete)  // 删除
	group.GET("detail", crud.Detail)  // 明细
}

type Model[T any] struct {
	Struct T        // 表对应的结构体
	DB     *gorm.DB // gorm
}

func (m *Model[T]) Create(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.BindJSON(&in); err != nil {
		log.Error("参数错误：", err)
		respx2.BuildException(ctx, respx2.ParamErr, err)
		return
	}
	err = m.DB.Create(&in).Error
	if err != nil {
		log.Error("对象新增失败 ： ", err)
		respx2.BuildError(ctx, err)
		return
	}
	respx2.BuildSuccess(ctx, nil)
}

func (m *Model[T]) Update(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.BindJSON(&in); err != nil {
		log.Error("参数错误：", err)
		respx2.BuildException(ctx, respx2.ParamErr, err)
		return
	}
	err = m.DB.Updates(&in).Error
	if err != nil {
		log.Error("对象更新失败 ： ", err)
		respx2.BuildError(ctx, err)
		return
	}
	respx2.BuildSuccess(ctx, nil)
}

func (m *Model[T]) Delete(ctx *gin.Context) {
	var err error
	var form modelx.IdString
	if err = ctx.ShouldBindQuery(&form); err != nil {
		log.Error("参数错误：", err)
		respx2.BuildException(ctx, respx2.ParamErr, err)
		return
	}
	err = m.DB.Where("id = '" + form.Id + "'").Delete(&m.Struct).Error
	if err != nil {
		log.Error("对象删除失败 ： ", err)
		respx2.BuildError(ctx, err)
		return
	}
	respx2.BuildSuccess(ctx, nil)
}

func (m *Model[T]) Detail(ctx *gin.Context) {
	var err error
	var form modelx.IdString
	if err = ctx.ShouldBindQuery(&form); err != nil {
		log.Error("参数错误：", err)
		respx2.BuildException(ctx, respx2.ParamErr, err)
		return
	}
	var result T
	err = m.DB.Where("id = '" + form.Id + "'").Find(&result).Error
	if err != nil {
		log.Error("对象查询失败 ： ", err)
		respx2.BuildError(ctx, err)
		return
	}
	respx2.BuildSuccess(ctx, result)
}
