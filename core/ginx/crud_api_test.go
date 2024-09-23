package ginx

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/core/gormx"
)

func TestCrudApiRouter(t *testing.T) {
	var ginEngine = gin.New()
	group := ginEngine.Group("/test")

	type User struct {
		Id   int64  `json:"id" gorm:"type:bigint; primary_key; comment:用户ID;"`
		Name string `json:"name" gorm:"type:varchar(100); not null; comment:姓名;"`
	}
	NewCrudApi[User](group, gormx.GetDB())
}
