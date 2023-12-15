package anyx

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/modelx"
	"github.com/go-xuan/quanx/common/respx"
	"github.com/go-xuan/quanx/runner/gormx"
)

func NewEChartApi[T any](router *gin.RouterGroup, relativePath string, source ...string) {
	var model = &Model[T]{DB: gormx.This().GetDB(source...)}
	group := router.Group(relativePath)
	group.POST("list", model.List)   // 列表
	group.POST("first", model.First) // 单插
	group.POST("pie", model.Pie)     // 饼状图
	group.POST("bar", model.Bar)     // 柱状图
	//group.POST("MultiBar", model.MultiBar) // 多项柱状图
}

func (m *Model[T]) Pie(ctx *gin.Context) {
	var err error
	var result []map[string]any
	var t T
	var lgdDataKey, dataKey = "lgd_data", "data"
	if method := reflect.ValueOf(&t).MethodByName("PieQuerySql"); method.IsValid() {
		querySql := method.Call([]reflect.Value{})[0].String()
		err = m.DB.Raw(querySql).Scan(&result).Error
		if err != nil {
			log.Error("查询失败 ： ", err)
			respx.BuildError(ctx, err)
			return
		}
	} else {
		err = m.DB.Model(&t).Scan(&result).Error
		if err != nil {
			log.Error("查询失败 ： ", err)
			respx.BuildError(ctx, err)
			return
		}
		var typeRef = reflect.TypeOf(t)
		for i := 0; i < typeRef.NumField(); i++ {
			if typeRef.Field(i).Tag.Get("echart") == lgdDataKey {
				lgdDataKey = strings.ToLower(typeRef.Field(i).Name)
			}
			if typeRef.Field(i).Tag.Get("echart") == dataKey {
				dataKey = strings.ToLower(typeRef.Field(i).Name)
			}
		}
	}

	var chart = modelx.PieChart{LgdData: make([]string, 0), Data: make([]any, 0)}
	for _, item := range result {
		if _, ok := item[lgdDataKey]; ok {
			chart.LgdData = append(chart.LgdData, item[lgdDataKey].(string))
			chart.Data = append(chart.Data, item[dataKey])
		}
	}
	respx.BuildSuccess(ctx, chart)
}

func (m *Model[T]) Bar(ctx *gin.Context) {
	var err error
	var result []map[string]any
	var t T
	var axisDataKey, dataKey = "axis_data", "data"
	if method := reflect.ValueOf(&t).MethodByName("BarQuerySql"); method.IsValid() {
		querySql := method.Call([]reflect.Value{})[0].String()
		err = m.DB.Raw(querySql).Scan(&result).Error
		if err != nil {
			log.Error("查询失败 ： ", err)
			respx.BuildError(ctx, err)
			return
		}
	} else {
		err = m.DB.Model(&t).Scan(&result).Error
		if err != nil {
			log.Error("查询失败 ： ", err)
			respx.BuildError(ctx, err)
			return
		}
		var typeRef = reflect.TypeOf(t)
		for i := 0; i < typeRef.NumField(); i++ {
			if typeRef.Field(i).Tag.Get("echart") == axisDataKey {
				axisDataKey = strings.ToLower(typeRef.Field(i).Name)
			}
			if typeRef.Field(i).Tag.Get("echart") == dataKey {
				dataKey = strings.ToLower(typeRef.Field(i).Name)
			}
		}
	}

	var chart = modelx.BarChart{AxisData: make([]string, 0), Data: make([]any, 0)}
	for _, item := range result {
		if _, ok := item[axisDataKey]; ok {
			chart.AxisData = append(chart.AxisData, item[axisDataKey].(string))
			chart.Data = append(chart.Data, item[dataKey])
		}
	}
	respx.BuildSuccess(ctx, chart)
}

func (m *Model[T]) MultiBar(ctx *gin.Context) {
	var err error
	var result []map[string]any
	var t T
	var axisDataKey, legendKey, dataKey = "axis_data", "legend", "data"
	if method := reflect.ValueOf(&t).MethodByName("BarQuerySql"); method.IsValid() {
		querySql := method.Call([]reflect.Value{})[0].String()
		err = m.DB.Raw(querySql).Scan(&result).Error
		if err != nil {
			log.Error("查询失败 ： ", err)
			respx.BuildError(ctx, err)
			return
		}
	} else {
		err = m.DB.Model(&t).Scan(&result).Error
		if err != nil {
			log.Error("查询失败 ： ", err)
			respx.BuildError(ctx, err)
			return
		}
		var typeRef = reflect.TypeOf(t)
		for i := 0; i < typeRef.NumField(); i++ {
			if typeRef.Field(i).Tag.Get("echart") == axisDataKey {
				axisDataKey = strings.ToLower(typeRef.Field(i).Name)
			}
			if typeRef.Field(i).Tag.Get("echart") == legendKey {
				legendKey = strings.ToLower(typeRef.Field(i).Name)
			}
			if typeRef.Field(i).Tag.Get("echart") == dataKey {
				dataKey = strings.ToLower(typeRef.Field(i).Name)
			}
		}
	}

	var chart = modelx.MultiBarChart{AxisData: make([]string, 0), Legend: make([]string, 0)}
	var axisDataMap = make(map[string]map[string]any)
	for _, item := range result {
		var legendMap = make(map[string]any)
		if axis, ok := item[axisDataKey]; ok {
			if _, ok = axisDataMap[axis.(string)]; ok {
				legendMap = axisDataMap[axis.(string)]
			}
		}
		legendMap[item[legendKey].(string)] = item[dataKey]
		axisDataMap[item[axisDataKey].(string)] = legendMap
	}
	respx.BuildSuccess(ctx, chart)
}
