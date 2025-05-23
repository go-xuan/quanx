package ginx

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/base/filex/excelx"
	"github.com/go-xuan/quanx/base/respx"
	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/common/modelx"
	"github.com/go-xuan/quanx/types/timex"
)

func NewCrudApi[T any](group *gin.RouterGroup, db *gorm.DB) {
	var api = &Model[T]{DB: db}
	group.GET("list", api.List)      // 列表
	group.POST("create", api.Create) // 新增
	group.POST("update", api.Update) // 修改
	group.GET("delete", api.Delete)  // 删除
	group.GET("detail", api.Detail)  // 明细
}

func NewExcelApi[T any](group *gin.RouterGroup, db *gorm.DB) {
	var api = &Model[T]{DB: db}
	group.POST("import", api.Import) // 新增
	group.POST("export", api.Export) // 修改
}

type Model[T any] struct {
	DB *gorm.DB
}

func (m *Model[T]) List(ctx *gin.Context) {
	var err error
	var result []*T
	if err = m.DB.Find(&result).Error; err != nil {
		respx.Error(ctx, err.Error())
	} else {
		respx.Success(ctx, result)
	}
}

func (m *Model[T]) Create(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.ShouldBindJSON(&in); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	if err = m.DB.Create(&in).Error; err != nil {
		respx.Error(ctx, err.Error())
	} else {
		respx.Success(ctx, nil)
	}
}

func (m *Model[T]) Update(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.ShouldBindJSON(&in); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	if err = m.DB.Updates(&in).Error; err != nil {
		respx.Error(ctx, err.Error())
	} else {
		respx.Success(ctx, nil)
	}
}

func (m *Model[T]) Delete(ctx *gin.Context) {
	var err error
	var id modelx.Id[string]
	if err = ctx.ShouldBindQuery(&id); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	var t T
	if err = m.DB.Where("id = ? ", id.Id).Delete(&t).Error; err != nil {
		respx.Error(ctx, err.Error())
	} else {
		respx.Success(ctx, nil)
	}
}

func (m *Model[T]) Detail(ctx *gin.Context) {
	var err error
	var id modelx.Id[string]
	if err = ctx.ShouldBindQuery(&id); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	var result T
	if err = m.DB.Where("id = ? ", id.Id).Find(&result).Error; err != nil {
		respx.Error(ctx, err.Error())
	} else {
		respx.Success(ctx, result)
	}
}

func (m *Model[T]) Import(ctx *gin.Context) {
	var err error
	var file modelx.File
	if err = ctx.ShouldBind(&file); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	var filePath = filepath.Join(constx.DefaultResourceDir, file.File.Filename)
	if err = ctx.SaveUploadedFile(file.File, filePath); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	var obj T
	var data []*T
	if data, err = excelx.ReadXlsxWithStruct(filePath, "", obj); err != nil {
		respx.Custom(ctx, http.StatusBadRequest, respx.NewResponseData(respx.ExportFailedCode, err.Error()))
		return
	}
	if err = m.DB.Model(obj).Create(&data).Error; err != nil {
		respx.Error(ctx, err.Error())
	} else {
		respx.Success(ctx, nil)
	}
}

func (m *Model[T]) Export(ctx *gin.Context) {
	var result []*T
	if err := m.DB.Find(&result).Error; err != nil {
		respx.Custom(ctx, http.StatusBadRequest, respx.NewResponseData(respx.ExportFailedCode, err.Error()))
		return
	}
	var filePath = filepath.Join(constx.DefaultResourceDir, time.Now().Format(timex.TimestampFmt)+".xlsx")
	if len(result) > 0 {
		if err := excelx.WriteXlsx(filePath, result); err != nil {
			respx.CustomError(ctx, respx.NewResponseData(respx.ExportFailedCode, err.Error()))
			return
		}
	}
	respx.File(ctx, filePath)
}
