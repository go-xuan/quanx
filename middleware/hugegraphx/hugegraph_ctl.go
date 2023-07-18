package hugegraphx

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/quanxiaoxuan/quanx/common/httpx"
)

var CTL *Control

// hugegraph制器
type Control struct {
	Config     *Config // hugegraph配置
	GremlinUrl string  // gremlin查询接口url
}

// 初始化redis控制器
func InitHugegraphCTL(conf *Config) {
	if CTL == nil {
		CTL = conf.NewHugegraphCTL()
		log.Info("初始化hugegraph连接-成功！", conf.Format())
	}
}

// gremlin查询API-get请求
func (ctl *Control) GremlinApiGet(gremlin string) (result Result, err error) {
	var bytes []byte
	bytes, err = httpx.GetHttp(ctl.GremlinUrl + `?gremlin=` + gremlin)
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
func (ctl *Control) GremlinApiPost(gremlin string) (result Result, err error) {
	var bytes []byte
	bytes, err = httpx.PostHttp(ctl.GremlinUrl, GremlinApiParam(gremlin))
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
func (ctl *Control) QueryVertexsPost(gremlin string) (vertexs Vertexs, requestId string, err error) {
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
func (ctl *Control) QueryVertexsGet(gremlin string) (vertexs Vertexs, requestId string, err error) {
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
func (ctl *Control) QueryEdgesPost(gremlin string) (edges Edges, requestId string, err error) {
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
func (ctl *Control) QueryEdgesGet(gremlin string) (edges Edges, requestId string, err error) {
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
func (ctl *Control) QueryPathsPost(gremlin string) (paths Paths, requestId string, err error) {
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
func (ctl *Control) QueryPathsGet(gremlin string) (paths Paths, requestId string, err error) {
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
func (ctl *Control) QueryValuesPost(gremlin string) (values []string, err error) {
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
func (ctl *Control) SchemaPostHttp(url string, param interface{}) (result interface{}, err error) {
	httpUrl := ctl.Config.SchemaHttpUrl(url)
	log.Info("请求URL : ", httpUrl)
	var resp []byte
	resp, err = httpx.PostHttp(httpUrl, param)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	return
}

// (顶点/边)附加属性
func (ctl *Control) SchemaPutHttp(url string, param interface{}) (result interface{}, err error) {
	httpUrl := ctl.Config.SchemaHttpUrl(url)
	log.Info("请求URL : ", httpUrl)
	var resp []byte
	resp, err = httpx.PutHttp(httpUrl, param)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	return
}

// (属性/顶点/边/索引)查询
func (ctl *Control) SchemaGetHttp(url string) (result interface{}, err error) {
	httpUrl := ctl.Config.SchemaHttpUrl(url)
	log.Info("请求URL : ", httpUrl)
	var resp []byte
	resp, err = httpx.GetHttp(httpUrl)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	return
}

// (属性/顶点/边/索引)删除
func (ctl *Control) SchemaDeleteHttp(url string) (result interface{}, err error) {
	httpUrl := ctl.Config.SchemaHttpUrl(url)
	log.Info("请求URL : ", httpUrl)
	var resp []byte
	resp, err = httpx.DeleteHttp(httpUrl)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	return
}
