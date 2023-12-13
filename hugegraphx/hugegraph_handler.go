package hugegraphx

import (
	"encoding/json"
	"github.com/go-xuan/quanx/httpx"
)

var handler *Handler

// hugegraph处理器
type Handler struct {
	Config     *Hugegraph // hugegraph配置
	GremlinUrl string     // Gremlin查询接口URL
	SchemaUrl  string     // schema操作接口URL
}

func This() *Handler {
	if handler == nil {
		panic("The gorm handler has not been initialized, please check the relevant config")
	}
	return handler
}

func (h *Handler) PropertykeysUrl() string {
	return h.SchemaUrl + Propertykeys
}

func (h *Handler) VertexlabelsUrl() string {
	return h.SchemaUrl + Vertexlabels
}

func (h *Handler) EdgelabelsUrl() string {
	return h.SchemaUrl + Edgelabels
}

func (h *Handler) IndexlabelsUrl() string {
	return h.SchemaUrl + Indexlabels
}

// gremlin查询API-get请求
func (h *Handler) GremlinGet(gremlin string, result interface{}) (requestId string, err error) {
	var bytes []byte
	if bytes, err = httpx.Get().Url(h.GremlinUrl + `?gremlin=` + gremlin).Do(); err != nil {
		return
	}
	var httpResult Result
	if err = json.Unmarshal(bytes, &httpResult); err != nil {
		return
	}
	requestId = httpResult.RequestId
	bytes, err = json.Marshal(httpResult.Result.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return
}

// gremlin查询API-Post请求
func (h *Handler) GremlinPost(gremlin string, result interface{}) (requestId string, err error) {
	var bindings interface{} // 构建绑定参数
	var aliases interface{}  // 构建图别名
	_ = json.Unmarshal([]byte(`{}`), &bindings)
	_ = json.Unmarshal([]byte(`{"graph": "hugegraph","g": "__g_hugegraph"}`), &aliases)
	var bytes []byte
	if bytes, err = httpx.Post().Url(h.GremlinUrl).Body(Param{
		Gremlin:  gremlin,
		Bindings: bindings,
		Language: "gremlin-groovy",
		Aliases:  aliases,
	}).Do(); err != nil {
		return
	}
	var httpResult Result
	if err = json.Unmarshal(bytes, &httpResult); err != nil {
		return
	}
	requestId = httpResult.RequestId
	bytes, err = json.Marshal(httpResult.Result.Data)
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
func (h *Handler) QueryVertexs(gremlin string) (vertexs Vertexs, requestId string, err error) {
	requestId, err = h.GremlinPost(gremlin, vertexs)
	if err != nil {
		return
	}
	return
}

// 查询边
func (h *Handler) QueryEdges(gremlin string) (edges Edges, requestId string, err error) {
	requestId, err = h.GremlinPost(gremlin, edges)
	if err != nil {
		return
	}
	return
}

// 查询path()
func (h *Handler) QueryPaths(gremlin string) (paths Paths, requestId string, err error) {
	requestId, err = h.GremlinPost(gremlin, paths)
	if err != nil {
		return
	}
	return
}

// 调用hugegraph的POST接口，返回属性值
func (h *Handler) QueryValues(gremlin string) (values []string, err error) {
	_, err = h.GremlinPost(gremlin, values)
	if err != nil {
		return
	}
	return
}
