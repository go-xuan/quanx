package execx

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/exec"
	"runtime"
)

// ExecCommand 执行命令
func ExecCommand(command string, in ...io.Reader) (string, string, error) {
	if runtime.GOOS == `windows` {
		return execCommandOnWindows(command, in...)
	} else {
		return execCommandOnLinux(command, in...)
	}
}

func execCommandOnLinux(command string, in ...io.Reader) (string, string, error) {
	cmd := exec.Command("/bin/bash", `-c`, command)
	cmd.Dir = "./"
	if len(in) > 0 {
		cmd.Stdin = ioutil.NopCloser(in[0])
	}
	return commandRun(cmd)
}

func execCommandOnWindows(command string, in ...io.Reader) (string, string, error) {
	cmd := exec.Command("cmd", `/C`, command)
	if len(in) > 0 {
		cmd.Stdin = ioutil.NopCloser(in[0])
	}
	return commandRun(cmd)
}

func commandRun(cmd *exec.Cmd) (string, string, error) {
	// 设置接收
	var stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	// 执行命令
	if err := cmd.Run(); err != nil {
		return stdout.String(), stderr.String(), err
	} else {
		return stdout.String(), stderr.String(), nil
	}
}
