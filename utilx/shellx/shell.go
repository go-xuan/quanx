package shellx

import (
	"bytes"
	"fmt"
	"os/exec"
)

// 执行shell命令
func ExecCommand(dir, shell string) (out string, err error) {
	var obf bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", shell)
	cmd.Dir = dir
	// 设置接收
	cmd.Stdout = &obf
	// 执行
	if err = cmd.Run(); err != nil {
		return
	}
	out = obf.String()
	fmt.Println(out)
	return
}
