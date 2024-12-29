package execx

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
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
	c.cmd.Stdin = ioutil.NopCloser(in)
	return c
}

func (c *Cmd) Run() (string, string, error) {
	if cmd := c.cmd; cmd != nil {
		var stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}
		cmd.Stdout, cmd.Stderr = stdout, stderr
		if err := cmd.Run(); err != nil {
			return stdout.String(), stderr.String(), err
		} else {
			return stdout.String(), stderr.String(), nil
		}
	}
	return "", "", errors.New("command not found")
}
