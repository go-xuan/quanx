package execx

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	out, _, err := ExecCommand("echo $GOPATH")
	fmt.Println(out)
	fmt.Println(err)
}
