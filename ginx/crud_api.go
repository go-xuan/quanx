package ginx

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/utilx/excelx"
	"github.com/go-xuan/utilx/timex"
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/quanx/modelx"
)

func NewCrudApi[T any](router *gin.RouterGroup, db *gorm.DB) {
	var api = &Model[T]{
		DB: db,
	}
	router.GET("list", api.List)      // 列表
	router.POST("create", api.Create) // 新增
	router.POST("update", api.Update) // 修改
	router.GET("delete", api.Delete)  // 删除
	router.GET("detail", api.Detail)  // 明细
}

func NewExcelApi[T any](group *gin.RouterGroup, db *gorm.DB) {
	var api = &Model[T]{DB: db}
	group.POST("import", api.Import) // 新增
	group.POST("export", api.Export) // 修改
}

type Model[T any] struct {
	DB *gorm.DB
}

// List 列表
func (m *Model[T]) List(ctx *gin.Context) {
	var err error
	var result []*T
	if err = m.DB.Find(&result).Error; err != nil {
		Error(ctx, err)
	} else {
		Success(ctx, result)
	}
}

// Create 新增
func (m *Model[T]) Create(ctx *gin.Context) {
	var err error
	var create T
	if err = ctx.ShouldBindJSON(&create); err != nil {
		ParamError(ctx, err)
		return
	}
	if err = m.DB.Create(&create).Error; err != nil {
		Error(ctx, err)
	} else {
		Success(ctx, nil)
	}
}

// Update 修改
func (m *Model[T]) Update(ctx *gin.Context) {
	var err error
	var update T
	if err = ctx.ShouldBindJSON(&update); err != nil {
		ParamError(ctx, err)
		return
	}
	if err = m.DB.Updates(&update).Error; err != nil {
		Error(ctx, err)
	} else {
		Success(ctx, nil)
	}
}

// Delete 删除
func (m *Model[T]) Delete(ctx *gin.Context) {
	var err error
	var id modelx.Id[string]
	if err = ctx.ShouldBindQuery(&id); err != nil {
		ParamError(ctx, err)
		return
	}
	var t T
	if err = m.DB.Where("id = ? ", id.Id).Delete(&t).Error; err != nil {
		Error(ctx, err)
	} else {
		Success(ctx, nil)
	}
}

// Detail 明细
func (m *Model[T]) Detail(ctx *gin.Context) {
	var err error
	var id modelx.Id[string]
	if err = ctx.ShouldBindQuery(&id); err != nil {
		ParamError(ctx, err)
		return
	}
	var result T
	if err = m.DB.Where("id = ? ", id.Id).Find(&result).Error; err != nil {
		Error(ctx, err)
	} else {
		Success(ctx, result)
	}
}

// Import 导入
func (m *Model[T]) Import(ctx *gin.Context) {
	var err error
	var file modelx.File
	if err = ctx.ShouldBind(&file); err != nil {
		ParamError(ctx, err)
		return
	}
	var filePath = filepath.Join(constx.DefaultResourceDir, file.File.Filename)
	if err = ctx.SaveUploadedFile(file.File, filePath); err != nil {
		ParamError(ctx, err)
		return
	}
	var obj T
	var data []*T
	if data, err = excelx.ReadAny(filePath, "", obj); err != nil {
		Custom(ctx, http.StatusBadRequest, NewResponse(ExportFailedCode, err.Error()))
		return
	}
	if err = m.DB.Model(obj).Create(&data).Error; err != nil {
		Error(ctx, err)
	} else {
		Success(ctx, nil)
	}
}

// Export 导出
func (m *Model[T]) Export(ctx *gin.Context) {
	var result []*T
	if err := m.DB.Find(&result).Error; err != nil {
		Custom(ctx, http.StatusBadRequest, NewResponse(ExportFailedCode, err.Error()))
		return
	}
	var filePath = filepath.Join(constx.DefaultResourceDir, time.Now().Format(timex.TimestampFmt)+".xlsx")
	if len(result) > 0 {
		if err := excelx.WriteAny(filePath, result); err != nil {
			CustomError(ctx, NewResponse(ExportFailedCode, err.Error()))
			return
		}
	}
	File(ctx, filePath)
}
