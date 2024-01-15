package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/commonx/modelx"
	"github.com/go-xuan/quanx/commonx/respx"
	"gorm.io/gorm"
)

type CrudModel interface {
	Path() string
	DB() *gorm.DB
}

type Crud interface {
	List(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Detail(ctx *gin.Context)
}

func NewCrudApi[T any](router *gin.RouterGroup, db *gorm.DB, path string) {
	var api Crud = &Model[T]{DB: db}
	group := router.Group(path)
	group.GET("list", api.List)      // 列表
	group.POST("create", api.Create) // 新增
	group.POST("update", api.Update) // 修改
	group.GET("delete", api.Delete)  // 删除
	group.GET("detail", api.Detail)  // 明细
}

type Model[T any] struct {
	Path string
	DB   *gorm.DB
}

func (m *Model[T]) List(ctx *gin.Context) {
	var err error
	var result []*T
	err = m.DB.Find(&result).Error
	respx.BuildResponse(ctx, result, err)
}

func (m *Model[T]) Create(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.BindJSON(&in); err != nil {
		respx.Exception(ctx, respx.ParamErr, err)
		return
	}
	err = m.DB.Create(&in).Error
	respx.BuildResponse(ctx, nil, err)
}

func (m *Model[T]) Update(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.BindJSON(&in); err != nil {
		respx.Exception(ctx, respx.ParamErr, err)
		return
	}
	err = m.DB.Updates(&in).Error
	respx.BuildResponse(ctx, nil, err)
}

func (m *Model[T]) Delete(ctx *gin.Context) {
	var err error
	var form modelx.Id[string]
	if err = ctx.ShouldBindQuery(&form); err != nil {
		respx.Exception(ctx, respx.ParamErr, err)
		return
	}
	var t T
	err = m.DB.Where("id = ? ", form.Id).Delete(&t).Error
	respx.BuildResponse(ctx, nil, err)
}

func (m *Model[T]) Detail(ctx *gin.Context) {
	var err error
	var form modelx.Id[string]
	if err = ctx.ShouldBindQuery(&form); err != nil {
		respx.Exception(ctx, respx.ParamErr, err)
		return
	}
	var result T
	err = m.DB.Where("id = ? ", form.Id).Find(&result).Error
	respx.BuildResponse(ctx, result, err)
}
