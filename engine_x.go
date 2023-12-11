package quanx

type X[T any] interface {
	Format() string
	Init()
}
