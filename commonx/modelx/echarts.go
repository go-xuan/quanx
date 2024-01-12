package modelx

// eCharts饼图（适用环形图）
type PieChart struct {
	LgdData []string `json:"lgdData"`
	Data    []any    `json:"data"`
}

// eCharts条形图
type BarChart struct {
	AxisData []string `json:"axisData"`
	Data     []any    `json:"data"`
}

// eCharts多项条形图
type MultiBarChart struct {
	AxisData []string `json:"axisData"`
	Legend   []string `json:"legend"`
	Data     [][]any  `json:"data"`
	Values   []any    `json:"values"`
}

// eCharts折线图
type LineChart struct {
	AxisData []string `json:"axisData"`
	Legend   []string `json:"legend"`
	Data     [][]any  `json:"data"`
}
