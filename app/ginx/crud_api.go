package ginx

import (
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/app/constx"
	"github.com/go-xuan/quanx/app/modelx"
	"github.com/go-xuan/quanx/net/respx"
	"github.com/go-xuan/quanx/os/filex/excelx"
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
	err = m.DB.Find(&result).Error
	respx.Ctx(ctx).Response(result, err)
}

func (m *Model[T]) Create(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.BindJSON(&in); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	err = m.DB.Create(&in).Error
	respx.Ctx(ctx).Response(nil, err)
}

func (m *Model[T]) Update(ctx *gin.Context) {
	var err error
	var in T
	if err = ctx.BindJSON(&in); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	err = m.DB.Updates(&in).Error
	respx.Ctx(ctx).Response(nil, err)
}

func (m *Model[T]) Delete(ctx *gin.Context) {
	var err error
	var form modelx.Id[string]
	if err = ctx.ShouldBindQuery(&form); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	var t T
	err = m.DB.Where("id = ? ", form.Id).Delete(&t).Error
	respx.Ctx(ctx).Response(nil, err)
}

func (m *Model[T]) Detail(ctx *gin.Context) {
	var err error
	var form modelx.Id[string]
	if err = ctx.ShouldBindQuery(&form); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	var result T
	err = m.DB.Where("id = ? ", form.Id).Find(&result).Error
	respx.Ctx(ctx).Response(nil, err)
}

func (m *Model[T]) Import(ctx *gin.Context) {
	var err error
	var form modelx.File
	if err = ctx.ShouldBind(&form); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	var filePath = filepath.Join(constx.DefaultResourceDir, form.File.Filename)
	if err = ctx.SaveUploadedFile(form.File, filePath); err != nil {
		respx.ParamError(ctx, err)
		return
	}
	var obj T
	var data []*T
	if data, err = excelx.ExcelReaderAny(filePath, "", obj); err != nil {
		respx.Error(ctx, respx.ImportFailedCode, err)
		return
	}
	err = m.DB.Model(obj).Create(&data).Error
	respx.Ctx(ctx).Response(nil, err)
}

func (m *Model[T]) Export(ctx *gin.Context) {
	var result []*T
	if err := m.DB.Find(&result).Error; err != nil {
		respx.Error(ctx, respx.ImportFailedCode, err)
		return
	}
	var filePath = filepath.Join(constx.DefaultResourceDir, time.Now().Format(timex.TimestampFmt)+".xlsx")
	var obj T
	if err := excelx.ExcelWriter(filePath, obj, result); err != nil {
		respx.Error(ctx, respx.ExportFailedCode, err)
	} else {
		respx.BuildExcelByFile(ctx, filePath)
	}
}
