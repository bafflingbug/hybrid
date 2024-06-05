package hybrid

// Warp 包装, 通过warp实现对nil的存储
type Warp[T any] struct {
	Value T `json:"value"`
}

func NewWarp[T any](val T) *Warp[T] {
	return &Warp[T]{
		Value: val,
	}
}

func NewWarpEmpty[T any]() *Warp[T] {
	return &Warp[T]{}
}
