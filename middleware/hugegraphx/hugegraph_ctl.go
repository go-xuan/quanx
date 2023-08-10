package hugegraphx

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/quanxiaoxuan/quanx/common/httpx"
)

var CTL *Controller

// hugegraph制器
type Controller struct {
	Config     *Config // hugegraph配置
	GremlinUrl string  // gremlin查询接口url
}

// 初始化redis控制器
func Init(conf *Config) {
	if CTL == nil {
		CTL = &Controller{Config: conf, GremlinUrl: conf.GremlinHttpUrl()}
		log.Info("初始化hugegraph连接-成功！", conf.Format())
	}
}

// gremlin查询API-get请求
func (ctl *Controller) GremlinApiGet(gremlin string) (result Result, err error) {
	var bytes []byte
	bytes, err = httpx.HttpGet(ctl.GremlinUrl + `?gremlin=` + gremlin)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return
}

// gremlin查询API-POST请求参数
func GremlinApiParam(gremlin string) Param {
	var bindings interface{} // 构建绑定参数
	var aliases interface{}  // 构建图别名
	_ = json.Unmarshal([]byte(`{}`), &bindings)
	_ = json.Unmarshal([]byte(`{"graph": "hugegraph","g": "__g_hugegraph"}`), &aliases)
	return Param{Gremlin: gremlin, Bindings: bindings, Language: "gremlin-groovy", Aliases: aliases}
}

// gremlin查询API-Post请求
func (ctl *Controller) GremlinApiPost(gremlin string) (result Result, err error) {
	var bytes []byte
	bytes, err = httpx.HttpPost(ctl.GremlinUrl, GremlinApiParam(gremlin))
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return
}

// 查询顶点
func (ctl *Controller) QueryVertexsPost(gremlin string) (vertexs Vertexs, requestId string, err error) {
	var apiResult Result
	apiResult, err = ctl.GremlinApiPost(gremlin)
	if err != nil {
		return
	}
	requestId = apiResult.RequestId
	// 将Result.Data序列化为json
	var bytes []byte
	bytes, err = json.Marshal(apiResult.Result.Data)
	if err != nil {
		return
	}
	// 将json反序列化为interface{}
	err = json.Unmarshal(bytes, &vertexs)
	if err != nil {
		return
	}
	return
}

// 查询顶点
func (ctl *Controller) QueryVertexsGet(gremlin string) (vertexs Vertexs, requestId string, err error) {
	var apiResult Result
	apiResult, err = ctl.GremlinApiGet(gremlin)
	if err != nil {
		return
	}
	requestId = apiResult.RequestId
	// 将Result.Data序列化为json
	var bytes []byte
	bytes, err = json.Marshal(apiResult.Result.Data)
	if err != nil {
		return
	}
	// 将json反序列化为interface{}
	err = json.Unmarshal(bytes, &vertexs)
	if err != nil {
		return
	}
	return
}

// 查询边
func (ctl *Controller) QueryEdgesPost(gremlin string) (edges Edges, requestId string, err error) {
	var apiResult Result
	apiResult, err = ctl.GremlinApiPost(gremlin)
	if err != nil {
		return
	}
	requestId = apiResult.RequestId
	// 将Result.Data序列化为json
	var bytes []byte
	bytes, err = json.Marshal(apiResult.Result.Data)
	if err != nil {
		return
	}
	// 将json反序列化
	err = json.Unmarshal(bytes, &edges)
	if err != nil {
		return
	}
	return
}

// 查询边
func (ctl *Controller) QueryEdgesGet(gremlin string) (edges Edges, requestId string, err error) {
	var apiResult Result
	apiResult, err = ctl.GremlinApiGet(gremlin)
	if err != nil {
		return
	}
	requestId = apiResult.RequestId
	// 将Result.Data序列化为json
	var bytes []byte
	bytes, err = json.Marshal(apiResult.Result.Data)
	if err != nil {
		return
	}
	// 将json反序列化
	err = json.Unmarshal(bytes, &edges)
	if err != nil {
		return
	}
	return
}

// 查询path()
func (ctl *Controller) QueryPathsPost(gremlin string) (paths Paths, requestId string, err error) {
	var apiResult Result
	apiResult, err = ctl.GremlinApiPost(gremlin)
	if err != nil {
		return
	}
	requestId = apiResult.RequestId
	// 将interface{}序列化为json
	var bytes []byte
	bytes, err = json.Marshal(apiResult.Result.Data)
	if err != nil {
		return
	}
	// 将json反序列化
	err = json.Unmarshal(bytes, &paths)
	return
}

// 查询path()
func (ctl *Controller) QueryPathsGet(gremlin string) (paths Paths, requestId string, err error) {
	var apiResult Result
	apiResult, err = ctl.GremlinApiGet(gremlin)
	if err != nil {
		return
	}
	requestId = apiResult.RequestId
	// 将interface{}序列化为json
	var bytes []byte
	bytes, err = json.Marshal(apiResult.Result.Data)
	if err != nil {
		return
	}
	// 将json反序列化
	err = json.Unmarshal(bytes, &paths)
	return
}

// 调用hugegraph的POST接口，返回属性值
func (ctl *Controller) QueryValuesPost(gremlin string) (values []string, err error) {
	var apiResult Result
	apiResult, err = ctl.GremlinApiPost(gremlin)
	if err != nil {
		return
	}
	var bytes []byte
	bytes, err = json.Marshal(apiResult.Result.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &values)
	if err != nil {
		return
	}
	return
}

// (属性/顶点/边/索引)新增
func (ctl *Controller) SchemaPostHttp(url string, param interface{}) (result interface{}, err error) {
	httpUrl := ctl.Config.SchemaHttpUrl(url)
	log.Info("请求URL : ", httpUrl)
	var resp []byte
	resp, err = httpx.HttpPost(httpUrl, param)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	return
}

// (顶点/边)附加属性
func (ctl *Controller) SchemaPutHttp(url string, param interface{}) (result interface{}, err error) {
	httpUrl := ctl.Config.SchemaHttpUrl(url)
	log.Info("请求URL : ", httpUrl)
	var resp []byte
	resp, err = httpx.HttpPut(httpUrl, param)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	return
}

// (属性/顶点/边/索引)查询
func (ctl *Controller) SchemaGetHttp(url string) (result interface{}, err error) {
	httpUrl := ctl.Config.SchemaHttpUrl(url)
	log.Info("请求URL : ", httpUrl)
	var resp []byte
	resp, err = httpx.HttpGet(httpUrl)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	return
}

// (属性/顶点/边/索引)删除
func (ctl *Controller) SchemaDeleteHttp(url string) (result interface{}, err error) {
	httpUrl := ctl.Config.SchemaHttpUrl(url)
	log.Info("请求URL : ", httpUrl)
	var resp []byte
	resp, err = httpx.HttpDelete(httpUrl)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	return
}
