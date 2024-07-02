package execx

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"runtime"
)

// ExecCommandAndLog 执行命令并记录日志
func ExecCommandAndLog(cmd string, msg string) (out string, err error) {
	if runtime.GOOS == `windows` {
		return ExecWindowsCommandAndLog(cmd, msg)
	} else {
		return ExecLinuxCommandAndLog(cmd, msg)
	}
}

// ExecLinuxCommandAndLog 执行linux命令并记录日志
func ExecLinuxCommandAndLog(cmd string, msg string) (out string, err error) {
	log.WithField(`cmd`, cmd).Info(msg)
	if out, err = ExecCommandOnLinux(cmd); err != nil {
		log.WithField(`cmd`, cmd).Error(err)
		return
	}
	return
}

// ExecWindowsCommandAndLog 执行windows命令并记录日志
func ExecWindowsCommandAndLog(cmd string, msg string) (out string, err error) {
	log.WithField(`cmd`, cmd).Info(msg)
	if out, err = ExecCommandOnWindows(cmd); err != nil {
		log.WithField(`cmd`, cmd).Error(err)
		return
	}
	return
}

func ExecCommandOnLinux(command string) (string, error) {
	return CommandRun(BuildCommand(`./`, `/bin/bash`, `-c`, command))
}

func ExecCommandOnWindows(command string) (string, error) {
	return CommandRun(BuildCommand(``, `cmd`, `/C`, command))
}

func CommandRun(cmd *exec.Cmd) (string, error) {
	// 设置接收
	var stdout, stderr = &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	// 执行命令
	if err := cmd.Run(); err != nil {
		return "\n" + stdout.String() + "\n" + stderr.String(), err
	} else {
		return "\n" + stdout.String(), nil
	}
}

func BuildCommand(dir string, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd
}
