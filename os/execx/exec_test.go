package execx

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	stdout, stderr, err := Command("echo $GOPATH").Run()
	fmt.Println("stdout:", stdout)
	fmt.Println("stderr:", stderr)
	fmt.Println(err)
}
