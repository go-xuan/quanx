package apix

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/commonx/modelx"
	"github.com/go-xuan/quanx/commonx/respx"
	"github.com/go-xuan/quanx/importx/gormx"
)

type Model[T any] struct {
	Struct T        // 表对应的结构体
	DB     *gorm.DB // gorm
}

func NewCrudApi[T any](router *gin.RouterGroup, relativePath string, source ...string) {
	var model = &Model[T]{DB: gormx.This().GetDB(source...)}
	group := router.Group(relativePath)
	group.GET("list", model.List)      // 列表
	group.POST("create", model.Create) // 新增
	group.POST("update", model.Update) // 修改
	group.GET("delete", model.Delete)  // 删除
	group.GET("detail", model.Detail)  // 明细
}

func (m *Model[T]) Create(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.BindJSON(&in); err != nil {
		log.Error("参数错误：", err)
		respx.Exception(ctx, respx.ParamErr, err)
		return
	}
	err = m.DB.Create(&in).Error
	if err != nil {
		log.Error("对象新增失败 ： ", err)
		respx.BuildError(ctx, err)
		return
	}
	respx.BuildSuccess(ctx, nil)
}

func (m *Model[T]) Update(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.BindJSON(&in); err != nil {
		log.Error("参数错误：", err)
		respx.Exception(ctx, respx.ParamErr, err)
		return
	}
	err = m.DB.Updates(&in).Error
	if err != nil {
		log.Error("对象更新失败 ： ", err)
		respx.BuildError(ctx, err)
		return
	}
	respx.BuildSuccess(ctx, nil)
}

func (m *Model[T]) Delete(ctx *gin.Context) {
	var err error
	var form modelx.IdString
	if err = ctx.ShouldBindQuery(&form); err != nil {
		log.Error("参数错误：", err)
		respx.Exception(ctx, respx.ParamErr, err)
		return
	}
	var t T
	err = m.DB.Where("id = '" + form.Id + "'").Delete(&t).Error
	if err != nil {
		log.Error("对象删除失败 ： ", err)
		respx.BuildError(ctx, err)
		return
	}
	respx.BuildSuccess(ctx, nil)
}

func (m *Model[T]) Detail(ctx *gin.Context) {
	var err error
	var form modelx.IdString
	if err = ctx.ShouldBindQuery(&form); err != nil {
		log.Error("参数错误：", err)
		respx.Exception(ctx, respx.ParamErr, err)
		return
	}
	var result T
	err = m.DB.Where("id = '" + form.Id + "'").Find(&result).Error
	if err != nil {
		log.Error("对象查询失败 ： ", err)
		respx.BuildError(ctx, err)
		return
	}
	respx.BuildSuccess(ctx, result)
}

func (m *Model[T]) List(ctx *gin.Context) {
	var err error
	var result []*T
	err = m.DB.Find(&result).Error
	if err != nil {
		log.Error("对象查询失败 ： ", err)
		respx.BuildError(ctx, err)
		return
	}
	respx.BuildSuccess(ctx, result)
}

func (m *Model[T]) First(ctx *gin.Context) {
	var err error
	var result T
	err = m.DB.First(&result).Error
	if err != nil {
		log.Error("对象查询失败 ： ", err)
		respx.BuildError(ctx, err)
		return
	}
	respx.BuildSuccess(ctx, result)
}
