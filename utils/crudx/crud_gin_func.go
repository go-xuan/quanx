package crudx

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/quanxiaoxuan/quanx/common/respx"
)

func AddCrudRouter[T any](ginRouter *gin.RouterGroup, path string, db *gorm.DB) {
	crud := ginRouter.Group(path)
	var model Model[T]
	model.DB = db
	crud.POST("add", model.AddHandler)       // 新增
	crud.POST("update", model.UpdateHandler) // 修改
	crud.GET("delete", model.DeleteHandler)  // 删除
	crud.GET("detail", model.DetailHandler)  // 明细
}

func (m *Model[T]) AddHandler(context *gin.Context) {
	var err error
	var add T
	if err = context.BindJSON(&add); err != nil {
		log.Error("参数错误：", err)
		respx.BuildExceptionResponse(context, respx.ParamErr, err)
		return
	}
	if err = m.DbAdd(add); err != nil {
		respx.BuildErrorResponse(context, err.Error())
	} else {
		respx.BuildSuccessResponse(context, nil)
	}
}

func (m *Model[T]) UpdateHandler(context *gin.Context) {
	var err error
	var update T
	if err = context.BindJSON(&update); err != nil {
		log.Error("参数错误：", err)
		respx.BuildExceptionResponse(context, respx.ParamErr, err)
		return
	}
	if err = m.DbUpdate(update); err != nil {
		respx.BuildErrorResponse(context, err.Error())
	} else {
		respx.BuildSuccessResponse(context, nil)
	}
}

func (m *Model[T]) DeleteHandler(context *gin.Context) {
	var err error
	var form struct {
		Id string `form:"id" binding:"required"`
	}
	if err = context.ShouldBindQuery(&form); err != nil {
		log.Error("参数错误：", err)
		respx.BuildExceptionResponse(context, respx.ParamErr, err)
		return
	}
	if err = m.DbDelete(form.Id); err != nil {
		respx.BuildErrorResponse(context, err.Error())
	} else {
		respx.BuildSuccessResponse(context, nil)
	}
}

func (m *Model[T]) DetailHandler(context *gin.Context) {
	var err error
	var form struct {
		Id string `form:"id" binding:"required"`
	}
	if err = context.ShouldBindQuery(&form); err != nil {
		log.Error("参数错误：", err)
		respx.BuildExceptionResponse(context, respx.ParamErr, err)
		return
	}
	var obj T
	if obj, err = m.DbDetail(form.Id); err != nil {
		respx.BuildErrorResponse(context, err.Error())
	} else {
		respx.BuildSuccessResponse(context, obj)
	}
}
