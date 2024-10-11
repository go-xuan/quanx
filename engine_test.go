package quanx

import (
	"testing"
)

func TestEngineRun(t *testing.T) {
	var e = NewEngine()
	// 服务启动
	e.RUN()
}
