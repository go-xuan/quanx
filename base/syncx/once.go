package syncx

import "sync"

func OnceDo(f func()) {
	var once sync.Once
	once.Do(f)
}
