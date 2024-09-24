package httpx

import (
	"fmt"
	"testing"
)

func TestHttpClient(t *testing.T) {
	if resp, err := Get("https://www.baidu.com").Do(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusOK())
	}
}
