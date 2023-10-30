package modelx

type PrimaryKey struct {
	Id int64 `form:"id" json:"id" binding:"required"`
}
