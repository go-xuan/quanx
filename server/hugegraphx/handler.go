package hugegraphx

import (
	"encoding/json"

	"github.com/go-xuan/quanx/net/httpx"
)

var handler *Handler

// Handler hugegraph处理器
type Handler struct {
	config     *Hugegraph // hugegraph配置
	gremlinUrl string     // gremlin查询接口URL
	schemaUrl  string     // schema操作接口URL
}

func This() *Handler {
	if handler == nil {
		panic("the hugegraph handler has not been initialized, please check the relevant config")
	}
	return handler
}

func (h *Handler) PropertykeysUrl() string {
	return h.schemaUrl + Propertykeys
}

func (h *Handler) VertexlabelsUrl() string {
	return h.schemaUrl + Vertexlabels
}

func (h *Handler) EdgelabelsUrl() string {
	return h.schemaUrl + Edgelabels
}

func (h *Handler) IndexlabelsUrl() string {
	return h.schemaUrl + Indexlabels
}

// GremlinGet gremlin查询API-get请求
func GremlinGet[T any](result T, gremlin string) (requestId string, err error) {
	var res *httpx.Response
	if res, err = httpx.Get(This().gremlinUrl + `?gremlin=` + gremlin).Do(); err != nil {
		return
	}
	var resp ApiResp[any]
	if err = res.Unmarshal(&resp); err != nil {
		return
	}
	requestId = resp.RequestId
	var bytes []byte
	if bytes, err = json.Marshal(resp.Result.Data); err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &result); err != nil {
		return
	}
	return
}

// GremlinPost gremlin查询API-Post请求
func GremlinPost[T any](result T, gremlin string) (requestId string, err error) {
	var bindings, aliases any // 构建绑定参数和图别名
	_ = json.Unmarshal([]byte(`{}`), &bindings)
	_ = json.Unmarshal([]byte(`{"graph": "hugegraph","g": "__g_hugegraph"}`), &aliases)
	var res *httpx.Response
	if res, err = httpx.Post(This().gremlinUrl).Body(Param{
		Gremlin:  gremlin,
		Bindings: bindings,
		Language: "gremlin-groovy",
		Aliases:  aliases,
	}).Do(); err != nil {
		return
	}
	var resp ApiResp[T]
	if err = res.Unmarshal(&resp); err != nil {
		return
	}
	requestId = resp.RequestId
	result = resp.Result.Data
	return
}

// QueryVertexs 查询顶点
func QueryVertexs[T any](gremlin string) (data Vertexs[T], requestId string, err error) {
	if requestId, err = GremlinPost(data, gremlin); err != nil {
		return
	}
	return
}

// QueryEdges 查询边
func QueryEdges[T any](gremlin string) (data Edges[T], requestId string, err error) {
	if requestId, err = GremlinPost(data, gremlin); err != nil {
		return
	}
	return
}

// QueryPaths 查询path()
func QueryPaths[T any](gremlin string) (data Paths[T], requestId string, err error) {
	if requestId, err = GremlinPost(data, gremlin); err != nil {
		return
	}
	return
}

// QueryValues 调用hugegraph的POST接口，返回属性值
func QueryValues(gremlin string) (data []string, err error) {
	if _, err = GremlinPost(data, gremlin); err != nil {
		return
	}
	return
}
