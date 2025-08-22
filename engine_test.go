package quanx

import (
	"testing"
)

func TestEngineRun(t *testing.T) {
	// 服务启动
	NewEngine(
		SetPort(9995),
		Debug(),
	).RUN(t.Context())
}
