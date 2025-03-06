package httpx

import (
	"fmt"
	"strconv"
	"testing"
)

func TestAsyncRequest(t *testing.T) {
	async := NewAsyncRequest(10)
	for i := 0; i < 101; i++ {
		async.Add(Get("http://localhost:3456/tools/demo1").Trace(strconv.Itoa(i)).Debug())
	}
	success, failed := async.Do()
	for _, response := range success {
		fmt.Println("success :", response.Err)
	}
	for _, response := range failed {
		fmt.Println("failed :", response.Err)
	}
}
