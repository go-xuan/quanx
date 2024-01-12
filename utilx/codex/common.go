package codex

// 字段配置
type Field struct {
	Name    string `json:"name"`    // 字段名
	Type    string `json:"type"`    // 字段类型（数据库类型）
	Comment string `json:"comment"` // 备注
	Default string `json:"default"` // 默认值
}
