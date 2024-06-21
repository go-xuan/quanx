package shellx

import (
	"bytes"
	"os/exec"
)

// ExecCommand 执行命令（path指定执行目录,缺省则在当前目录执行）
func ExecCommand(command string, path ...string) (string, error) {
	var stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}
	cmd := exec.Command("/bin/bash", "-c", command)
	// 设置接收
	cmd.Stdout, cmd.Stderr = stdout, stderr
	// 设置执行路径
	if len(path) > 0 {
		cmd.Dir = path[0]
	} else {
		cmd.Dir = "./"
	}
	// 执行命令
	if err := cmd.Run(); err != nil {
		return stderr.String(), err
	} else {
		return stdout.String(), nil
	}
}

func Pwd(path string) string {

}
