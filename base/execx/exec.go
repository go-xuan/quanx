package execx

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"runtime"
)

func Command(command string) *Cmd {
	if runtime.GOOS == `windows` {
		return &Cmd{exec.Command("cmd", `/C`, command)}
	} else {
		return &Cmd{exec.Command("/bin/bash", `-c`, command)}
	}
}

type Cmd struct {
	cmd *exec.Cmd
}

func (c *Cmd) Dir(dir string) *Cmd {
	c.cmd.Dir = dir
	return c
}

func (c *Cmd) Stdin(in io.Reader) *Cmd {
	c.cmd.Stdin = io.NopCloser(in)
	return c
}

func (c *Cmd) Run() (string, string, error) {
	if c.cmd == nil {
		return "", "", errors.New("command instance is nil") // 更合适的错误信息
	}

	var stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}
	c.cmd.Stdout, c.cmd.Stderr = stdout, stderr
	err := c.cmd.Run()
	return stdout.String(), stderr.String(), err
}
