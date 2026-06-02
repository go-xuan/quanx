package ginx

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/utilx/excelx"
	"github.com/go-xuan/utilx/idx"
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/dbx"
	"github.com/go-xuan/quanx/modelx"
)

// DBGetter 数据库获取接口，用于解耦 ginx 与 dbx
type DBGetter interface {
	GetDB(source string) *gorm.DB
}

var dbGetter DBGetter

// SetDBGetter 设置数据库获取器
func SetDBGetter(getter DBGetter) {
	dbGetter = getter
}

// BindCrudRouter 新增crud路由
func BindCrudRouter[T any](router *gin.RouterGroup, source string) {
	api := &Model[T]{Source: source}
	router.GET("list", api.List)        // 列表
	router.GET("detail", api.Detail)    // 明细
	router.POST("create", api.Create)   // 新增
	router.PUT("update", api.Update)    // 修改
	router.DELETE("delete", api.Delete) // 删除
}

// BindExcelRouter 新增 Excel 相关路由
func BindExcelRouter[T any](group *gin.RouterGroup, source string) {
	api := &Model[T]{Source: source}
	group.POST("import", api.Import) // 导入
	group.POST("export", api.Export) // 导出
}

// Model 通用模型
type Model[T any] struct {
	Source string
	DB     *gorm.DB
}

// GetDB 获取数据库连接
func (m *Model[T]) GetDB() *gorm.DB {
	if m.DB == nil {
		if dbGetter != nil {
			m.DB = dbGetter.GetDB(m.Source)
		} else {
			m.DB = dbx.GetGormDB(m.Source)
		}
	}
	return m.DB
}

// List 列表（支持分页）
func (m *Model[T]) List(ctx *gin.Context) {
	var query modelx.Query
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ParamError(ctx, err)
		return
	}

	db := m.GetDB().Model(new(T))

	// 总数查询
	var total int64
	if err := db.Count(&total).Error; err != nil {
		Error(ctx, err)
		return
	}

	// 分页查询
	var result []*T
	if err := dbx.QueryPage(db, &query).Find(&result).Error; err != nil {
		Error(ctx, err)
		return
	}
	Success(ctx, query.BuildResp(result, total))
}

// Create 新增
func (m *Model[T]) Create(ctx *gin.Context) {
	var create T
	if err := ctx.ShouldBindJSON(&create); err != nil {
		ParamError(ctx, err)
		return
	}
	if err := m.GetDB().Create(&create).Error; err != nil {
		Error(ctx, err)
		return
	}
	Success(ctx, nil)
}

// Update 修改
func (m *Model[T]) Update(ctx *gin.Context) {
	var id modelx.Id[string]
	if err := ctx.ShouldBindQuery(&id); err != nil {
		ParamError(ctx, err)
		return
	}
	var update T
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ParamError(ctx, err)
		return
	}
	if err := m.GetDB().Where("id = ?", id.Id).Updates(&update).Error; err != nil {
		Error(ctx, err)
		return
	}
	Success(ctx, nil)
}

// Delete 删除
func (m *Model[T]) Delete(ctx *gin.Context) {
	var id modelx.Id[string]
	if err := ctx.ShouldBindQuery(&id); err != nil {
		ParamError(ctx, err)
		return
	}
	var t T
	if err := m.GetDB().Where("id = ? ", id.Id).Delete(&t).Error; err != nil {
		Error(ctx, err)
		return
	}
	Success(ctx, nil)
}

// Detail 明细
func (m *Model[T]) Detail(ctx *gin.Context) {
	var id modelx.Id[string]
	if err := ctx.ShouldBindQuery(&id); err != nil {
		ParamError(ctx, err)
		return
	}
	var result T
	if err := m.GetDB().Where("id = ? ", id.Id).Find(&result).Error; err != nil {
		Error(ctx, err)
		return
	}
	Success(ctx, result)
}

// Import 导入
func (m *Model[T]) Import(ctx *gin.Context) {
	var file modelx.File
	if err := ctx.ShouldBind(&file); err != nil {
		ParamError(ctx, err)
		return
	}
	path := filepath.Join("import", file.File.Filename)
	if err := ctx.SaveUploadedFile(file.File, path); err != nil {
		ParamError(ctx, err)
		return
	}
	var t T
	if data, err := excelx.ReadAny(path, "", t); err != nil {
		CustomResponse(ctx, NewResponse(ExportFailedCode, err.Error()))
		return
	} else if err = m.GetDB().Model(t).Create(&data).Error; err != nil {
		Error(ctx, err)
		return
	}
	Success(ctx, nil)
}

// Export 导出
func (m *Model[T]) Export(ctx *gin.Context) {
	var result []*T
	if err := m.GetDB().Find(&result).Error; err != nil {
		CustomResponse(ctx, NewResponse(ExportFailedCode, err.Error()))
		return
	}

	filePath := filepath.Join("export", idx.Timestamp()+".xlsx")
	if len(result) > 0 {
		if err := excelx.WriteAny(filePath, result); err != nil {
			CustomResponse(ctx, NewResponse(ExportFailedCode, err.Error()))
			return
		}
	}
	ctx.File(filePath)
}
