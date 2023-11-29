package codex

const StructTemplate = `
type {modelName}Query struct {
	Keyword string 				// 关键字
	Page *request.Page 	// 分页查询参数
}
`

const ControllerTemplate = `
import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func {modelName}Page(context *gin.Context) {
	var err error
	var in model.{modelName}Query
	if err = context.BindJSON(&in); err != nil {
		log.Error("参数错误：", err)
		respx.BuildException(context, respx.ParamErr, err)
		return
	}
	var result *respx.PageResponse
	result, err = service.{modelName}Page(&in)
	respx.BuildNormal(ctx, result, err)
}


func {modelName}List(context *gin.Context) {
	var err error
	var in model.{modelName}Query
	if err = context.BindJSON(&in); err != nil {
		log.Error("参数错误：", err)
		respx.BuildException(context, respx.ParamErr, err)
		return
	}
	var result []*table.{modelName}
	result, err = service.{modelName}List(&in)
	respx.BuildNormal(ctx, result, err)
}

func {modelName}Create(context *gin.Context) {
	var err error
	var in table.{modelName}
	if err = context.BindJSON(&in); err != nil {
		log.Error("参数错误：", err)
		respx.BuildException(context, respx.ParamErr, err)
		return
	}
	var result string
	result, err = service.{modelName}Create(&in)
	respx.BuildNormal(ctx, result, err)
}

func {modelName}Detail(context *gin.Context) {
	var err error
	var form modelx.IdString
	if err = context.ShouldBindQuery(&form); err != nil {
		log.Error("参数错误：", err)
		respx.BuildException(context, respx.ParamErr, err)
		return
	}
	var result *table.{modelName}
	result, err = service.{modelName}Detail(form.Id)
	respx.BuildNormal(ctx, result, err)
}
`

const LogicTemplate = `
func {modelName}Page(in *model.{modelName}Query) (resp *respx.PageResponse, err error) {
	var rows []*table.{modelName}
	var total int64
	rows, total, err = dao.{modelName}Page(in)
	if err != nil {
		return
	}
	resp = respx.BuildPageResp(in.Page.Page, rows, total)
	return
}

func {modelName}List(in *model.{modelName}Query) (result []*table.{modelName}, err error) {
	result, err = dao.{modelName}List(param)
	if err != nil {
		return
	}
	return
}

func {modelName}Create(in *table.{modelName}) (id string, err error) {
	in.Id = idx.SnowFlake().NewString()
	err = dao.{modelName}Create(in)
	if err != nil {
		return
	}
	id = in.Id
	return
}

func {modelName}Detail(id string) (result *table.{modelName}, err error) {
	result, err = dao.{modelName}Detail(id)
	if err != nil {
		return
	}
	return
}
`

const DaoTemplate = `
func {modelName}Page(query model.{modelName}Query) (rows []*table.{modelName}, total int64, err error) {
	selectSql := strings.Builder{}
	countSql := strings.Builder{}
	err = gormx.CTL.DB.Raw(selectSql.String()).Scan(&rows).Error
	if err != nil {
		return
	}
	err = gormx.CTL.DB.Raw(countSql.String()).Scan(&total).Error
	if err != nil {
		return 
	}
	return
}

func {modelName}List(in *model.{modelName}Query) (result []*table.{modelName}, err error) {
	sql := strings.Builder{}
	err = gormx.CTL.DB.Raw(sql.String()).Scan(&result).Error
	if err != nil {
		return
	}
	return
}

func {modelName}Create(in *table.{modelName}) (err error) {
	err = gormx.CTL.DB.Create(in).Error
	if err != nil {
		return
	}
	return
}

func {modelName}Detail(id string) (result *table.{modelName}, err error) {
	result = &table.{modelName}{}
	err = gormx.CTL.DB.First(result, id).Error
	if err != nil {
		return
	}
	if result.Id == "" {
		err = errors.New("此记录不存在")
		return
	}
	return
}
`
