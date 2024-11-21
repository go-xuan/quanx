package flagx

import (
	"fmt"
	"testing"
)

func TestCommand(t *testing.T) {
	command := NewCommand("test", "test command",
		StringOption("aaa", "option aaa", "aaa"),
		BoolOption("bbb", "option bbb", false),
	).SetExecutor(func() error {
		fmt.Println("execute command")
		return nil
	})

	// 注册命令
	command.Register()
	if err := Execute(); err != nil {
		fmt.Println(err)
	}
}
