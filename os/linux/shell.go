package linux

import (
	"bytes"
	"fmt"
	"os/exec"
)

// 执行shell命令
func ExecCommand(cmd string, dir ...string) (out string, err error) {
	var buffer bytes.Buffer
	command := exec.Command("/bin/bash", "-c", cmd)
	// 设置执行路径
	if len(dir) > 0 {
		command.Dir = dir[0]
	} else {
		command.Dir = "./"
	}
	// 设置接收
	command.Stdout = &buffer
	// 执行命令
	if err = command.Run(); err != nil {
		return
	}
	out = buffer.String()
	// 打印输出结果
	fmt.Println(out)
	return
}
