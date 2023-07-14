package shellx

import (
	"bytes"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// 执行shell命令
func ExecCommand(execPath, execShell string) (string, error) {
	fmt.Println(execShell)
	var outInfo bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", execShell)
	cmd.Dir = execPath
	// 设置接收
	cmd.Stdout = &outInfo
	// 执行
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return ``, err
	}
	fmt.Println(outInfo.String())
	return outInfo.String(), nil
}
