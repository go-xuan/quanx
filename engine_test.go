package quanx

import (
	"testing"
)

func TestEngineRun(t *testing.T) {
	// 服务启动
	NewEngine(Debug).RUN()
}
