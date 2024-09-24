package execx

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	out, err := ExecCommand("echo $GOPATH")
	fmt.Println(out)
	fmt.Println(err)
}
