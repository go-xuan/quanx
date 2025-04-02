package hugegraphx

import (
	"encoding/json"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/httpx"
)

var _handler *Client

func this() *Client {
	if _handler == nil {
		panic("hugegraph client not initialized, please check the relevant config")
	}
	return _handler
}

// Client hugegraph处理器
type Client struct {
	config     *Config // hugegraph配置
	gremlinUrl string  // gremlin查询接口URL
	schemaUrl  string  // schema操作接口URL
}

func (c *Client) Config() *Config {
	return c.config
}

func (c *Client) PropertykeysUrl() string {
	return c.schemaUrl + Propertykeys
}

func (c *Client) VertexlabelsUrl() string {
	return c.schemaUrl + Vertexlabels
}

func (c *Client) EdgelabelsUrl() string {
	return c.schemaUrl + Edgelabels
}

func (c *Client) IndexlabelsUrl() string {
	return c.schemaUrl + Indexlabels
}

func GetConfig() *Config {
	return this().Config()
}

// GremlinGet gremlin查询API-get请求
func GremlinGet[T any](result T, gremlin string) (string, error) {
	res, err := httpx.Get(GetConfig().GremlinUrl() + `?gremlin=` + gremlin).Do()
	if err != nil {
		return "", errorx.Wrap(err, "do gremlin query failed")
	}
	var resp ApiResp[any]
	if err = res.Unmarshal(&resp); err != nil {
		return "", errorx.Wrap(err, "unmarshal gremlin response failed")
	}
	requestId := resp.RequestId
	var bytes []byte
	if bytes, err = json.Marshal(resp.Result.Data); err != nil {
		return "", errorx.Wrap(err, "marshal result failed")
	}
	if err = json.Unmarshal(bytes, &result); err != nil {
		return "", errorx.Wrap(err, "unmarshal result failed")
	}
	return requestId, nil
}

// GremlinPost gremlin查询API-Post请求
func GremlinPost[T any](result T, gremlin string) (string, error) {
	var bindings, aliases any // 构建绑定参数和图别名
	_ = json.Unmarshal([]byte(`{}`), &bindings)
	_ = json.Unmarshal([]byte(`{"graph": "hugegraph","g": "__g_hugegraph"}`), &aliases)
	res, err := httpx.Post(GetConfig().GremlinUrl()).Body(Param{
		Gremlin:  gremlin,
		Bindings: bindings,
		Language: "gremlin-groovy",
		Aliases:  aliases,
	}).Do()
	if err != nil {
		return "", errorx.Wrap(err, "do gremlin query failed")
	}
	var resp ApiResp[T]
	if err = res.Unmarshal(&resp); err != nil {
		return "", errorx.Wrap(err, "unmarshal gremlin response failed")
	}
	requestId := resp.RequestId
	result = resp.Result.Data
	return requestId, nil
}

// QueryVertexs 查询顶点
func QueryVertexs[T any](gremlin string) (Vertexs[T], string, error) {
	var data Vertexs[T]
	if requestId, err := GremlinPost(data, gremlin); err != nil {
		return data, "", errorx.Wrap(err, "gremlin query failed")
	} else {
		return data, requestId, nil
	}
}

// QueryEdges 查询边
func QueryEdges[T any](gremlin string) (Edges[T], string, error) {
	var data Edges[T]
	if requestId, err := GremlinPost(data, gremlin); err != nil {
		return data, "", errorx.Wrap(err, "gremlin query failed")
	} else {
		return data, requestId, nil
	}
}

// QueryPaths 查询path()
func QueryPaths[T any](gremlin string) (Paths[T], string, error) {
	var data Paths[T]
	if requestId, err := GremlinPost(data, gremlin); err != nil {
		return data, "", errorx.Wrap(err, "gremlin query failed")
	} else {
		return data, requestId, nil
	}
}

// QueryValues 调用hugegraph的POST接口，返回属性值
func QueryValues(gremlin string) ([]string, error) {
	var data []string
	if _, err := GremlinPost(data, gremlin); err != nil {
		return data, errorx.Wrap(err, "gremlin query failed")
	}
	return data, nil
}
