package httpx

import (
	"fmt"
	"testing"
)

func TestRequest(t *testing.T) {
	if resp, err := Get("http://localhost:3456/tools/demo").Debug().Do(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusOK())
	}
}
