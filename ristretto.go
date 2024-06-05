package hybrid

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
)

type ristrettoCache[T any] struct {
	cache *ristretto.Cache

	cost int64
	ttl  time.Duration
	opt  *Option
}

// Set 实现接口
func (r *ristrettoCache[T]) Set(ctx context.Context, key string, val T, isEmpty bool) error {
	ttl := r.ttl
	if isEmpty {
		if !r.opt.Empty {
			return nil
		}
		if r.opt.EmptyTtl > 0 {
			ttl = r.opt.EmptyTtl
		}
	}
	if !r.cache.SetWithTTL(r.key(key), val, r.cost, ttl) {
		return fmt.Errorf("ristretto could not set key %s", key)
	}
	return nil
}

// Get 实现接口
func (r *ristrettoCache[T]) Get(ctx context.Context, key string) (T, error) {
	v, ok := r.cache.Get(r.key(key))
	if !ok {
		return *new(T), NotFindCache
	}
	if t, ok2 := v.(T); ok2 {
		return t, nil
	}
	t := *new(T)
	return t, fmt.Errorf("ristretto find key %s but [%T] is not [%T]", key, v, t)
}

// Del 实现接口
func (r *ristrettoCache[T]) Del(ctx context.Context, key string) error {
	r.cache.Del(r.key(key))
	return nil
}

func (r *ristrettoCache[T]) key(key string) string {
	if r.opt.Prefix != "" {
		return fmt.Sprintf("[%s]%s", r.opt.Prefix, key)
	}
	return key
}

// WithRistretto 使用ristretto本地缓存
func WithRistretto[T any](cache *ristretto.Cache, cost int64, ttl time.Duration, opts ...Options) Cache[T] {
	opt := &Option{}
	for _, optFn := range opts {
		optFn(opt)
	}
	return &ristrettoCache[T]{
		cache: cache,
		cost:  cost,
		ttl:   ttl,
		opt:   opt,
	}
}
