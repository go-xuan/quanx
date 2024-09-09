package idx

import (
	"time"

	"github.com/google/uuid"
)

func UUID() string {
	return uuid.NewString()
}

func TimeUnix() int64 {
	return time.Now().UnixMilli()
}

func Timestamp() string {
	return time.Now().UTC().Format("20060102150405")
}
