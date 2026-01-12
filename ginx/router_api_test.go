package ginx

import (
	"testing"

	"github.com/go-xuan/quanx/constx"
)

func TestCrudApiRouter(t *testing.T) {
	e := DefaultEngine()
	SetDebugMode(e)
	type User struct {
		Id   int64  `json:"id" gorm:"type:bigint; primary_key; comment:用户ID;"`
		Name string `json:"name" gorm:"type:varchar(100); not null; comment:姓名;"`
	}
	BindCrudRouter[User](e.Group("/test/user"), constx.DefaultSource)
}
