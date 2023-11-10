package modelx

type IdInt64 struct {
	Id int64 `form:"id" json:"id" binding:"required"`
}

type IdsInt64 struct {
	Ids []int64 `form:"ids" json:"ids" binding:"required"`
}

type IdString struct {
	Id string `form:"id" json:"id" binding:"required"`
}

type IdsString struct {
	Ids string `form:"ids" json:"ids" binding:"required"`
}
