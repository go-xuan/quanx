package httpx

import (
	"sync"
)

// NewAsyncRequest 初始化异步请求器
func NewAsyncRequest(concurrency int) *AsyncRequester {
	return &AsyncRequester{
		concurrency: concurrency,
	}
}

// AsyncRequester 异步请求器
type AsyncRequester struct {
	requests    []*Request // 全部请求
	concurrency int        // 并发数
}

// Add 添加请求
func (a *AsyncRequester) Add(req ...*Request) *AsyncRequester {
	a.requests = append(a.requests, req...)
	return a
}

// Do 执行请求
func (a *AsyncRequester) Do() (success, failed []*AsyncResponse) {
	// 初始化管道和wg
	var total = len(a.requests)
	if a.concurrency < 1 {
		a.concurrency = 10 // 默认10个并发
	}
	var chanBuffer = total / a.concurrency // 用于接收的并发管道缓冲
	if r := total % a.concurrency; r > 0 {
		chanBuffer = chanBuffer + 1
	}
	var reqChan = make(chan *Request, total)
	var respChan = make(chan *AsyncResponse, chanBuffer)
	var wg = &sync.WaitGroup{}

	// 将所有请求添加进管道
	for _, request := range a.requests {
		reqChan <- request
	}
	close(reqChan)

	// 并发处理异步请求
	for i := 0; i < a.concurrency; i++ {
		wg.Add(1)
		go doAsync(reqChan, respChan, wg)
	}

	go func() {
		wg.Wait()
		close(respChan)
	}()

	// 从管道获取请求结果
	for resp := range respChan {
		if resp.IsSuccess() {
			success = append(success, resp)
		} else {
			failed = append(failed, resp)
		}
	}
	return success, failed
}

// AsyncResponse 异步请求响应结果
type AsyncResponse struct {
	Resp *Response
	Err  error
}

func (a *AsyncResponse) IsSuccess() bool {
	if a.Err == nil && a.Resp != nil {
		return a.Resp.StatusOK()
	}
	return false
}

// 异步请求
func doAsync(reqChan <-chan *Request, respChan chan<- *AsyncResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	for request := range reqChan {
		if resp, err := request.Do(); err != nil {
			respChan <- &AsyncResponse{Err: err}
		} else {
			respChan <- &AsyncResponse{Resp: resp}
		}
	}
}
