package modelx

import "time"

type Base interface {
	GetCreateUser() string
	GetCreateTime() time.Time
	GetUpdateUser() string
	GetUpdateTime() time.Time
}
