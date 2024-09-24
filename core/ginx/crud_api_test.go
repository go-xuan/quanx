package ginx

import (
	"gorm.io/gorm"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCrudApiRouter(t *testing.T) {
	type User struct {
		Id   int64  `json:"id" gorm:"type:bigint; primary_key; comment:用户ID;"`
		Name string `json:"name" gorm:"type:varchar(100); not null; comment:姓名;"`
	}
	NewCrudApi[User](
		gin.New().Group("/test/user"),
		&gorm.DB{})
}
