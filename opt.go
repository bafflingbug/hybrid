package hybrid

import "time"

// Option 选项
type Option struct {
	Prefix   string
	Empty    bool
	EmptyTtl time.Duration
}

// Options 选项
type Options func(option *Option)

// WithPrefix 写db添加前缀
func WithPrefix(prefix string) Options {
	return func(option *Option) {
		option.Prefix = prefix
	}
}

// WithCacheEmpty 是否写入空数据
func WithCacheEmpty(empty bool) Options {
	return func(option *Option) {
		option.Empty = empty
	}
}

// WithCacheEmptyTtl 空数据ttl
func WithCacheEmptyTtl(emptyTtl time.Duration) Options {
	return func(option *Option) {
		option.EmptyTtl = emptyTtl
	}
}
