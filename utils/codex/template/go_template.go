package template

const GoParam = `
type {modelName}Query struct {
	SearchKey string 				// 关键字
	PageParam *request.PageParam 	// 分页查询参数
}
`

const GoController = `
import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func {modelName}QueryHandler(context *gin.Context) {
	var err error
	var form struct {
		Id string
	}
	if err = context.ShouldBindQuery(&form); err != nil {
		log.Error("参数错误：", err)
		response.BuildExceptionResponse(context, response.ParamError, err)
		return
	}
	var result entity.{modelName}
	if result, err = service.{modelName}Query(form.Id); err != nil {
		response.BuildExceptionResponse(context, response.ServerError, err)
	} else {
		response.BuildSuccessResponse(context, result)
	}
}

func {modelName}AddHandler(context *gin.Context) {
	var err error
	var param entity.{modelName}
	if err = context.BindJSON(&param); err != nil {
		log.Error("参数错误：", err)
		response.BuildExceptionResponse(context, response.ParamError, err)
		return
	}
	if err = service.{modelName}Add(param); err != nil {
		response.BuildExceptionResponse(context, response.ServerError, err)
	} else {
		response.BuildSuccessResponse(context, nil)
	}
}

func {modelName}ListHandler(context *gin.Context) {
	var err error
	var param params.{modelName}Query
	if err = context.BindJSON(&param); err != nil {
		log.Error("参数错误：", err)
		response.BuildExceptionResponse(context, response.ParamError, err)
		return
	}
	var result entity.{modelName}List
	if result, err = service.{modelName}List(param); err != nil {
		response.BuildExceptionResponse(context, response.ServerError, err)
	} else {
		response.BuildSuccessResponse(context, result)
	}
}

func {modelName}PageHandler(context *gin.Context) {
	var err error
	var param params.{modelName}Query
	if err = context.BindJSON(&param); err != nil {
		log.Error("参数错误：", err)
		response.BuildExceptionResponse(context, response.ParamError, err)
		return
	}
	var result *response.PageResponse
	if result, err = service.{modelName}Page(param); err != nil {
		response.BuildExceptionResponse(context, response.ServerError, err)
	} else {
		response.BuildSuccessResponse(context, result)
	}
}
`

const GoService = `
func {modelName}Query(id string) (result entity.{modelName}, err error) {
	return dao.{modelName}Query(id)
}

func {modelName}Add(param entity.{modelName}) (err error) {
	return dao.{modelName}Add(param)
}

func {modelName}List(param params.{modelName}Query) (result entity.{modelName}List, err error) {
	return dao.{modelName}List(param)
}

func {modelName}Page(param params.{modelName}Query) (*response.PageResponse, error) {
	var resultList entity.{modelName}List
	var total int64
	var err error
	resultList, total, err = dao.{modelName}Page(param)
	if err != nil {
		return nil, err
	}
	return response.BuildPageData(param.PageParam, resultList, total), nil
}
`

const GoDao = `
func {modelName}Query(id string) (result entity.{modelName}, err error) {
	err = database.DB.First(&result, id).Error
	return
}

func {modelName}Add(param entity.{modelName}) (err error) {
	err = database.DB.Create(&param).Error
	return
}

func {modelName}List(param params.{modelName}Query) (resultList entity.{modelName}List, err error) {
	sql := strings.Builder{}
	err = database.DB.Raw(sql.String()).Scan(&resultList).Error
	return
}

func {modelName}Page(query params.{modelName}Query) (resultList entity.{modelName}List, total int64, err error) {
	selectSql := strings.Builder{}
	countSql := strings.Builder{}
	err = database.DB.Raw(selectSql.String()).Scan(&resultList).Error
	if err != nil {
		return
	}
	err = database.DB.Raw(countSql.String()).Scan(&total).Error
	if err != nil {
		return 
	}
	return
}
`
