// Package hybrid 多级缓存
package hybrid

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type Cache[T any] interface {
	Set(ctx context.Context, key string, val T, isEmpty bool) error
	Get(ctx context.Context, key string) (T, error)
	Del(ctx context.Context, key string) error
}

type layer[T any] struct {
	cache Cache[T]
	next  *layer[T]

	unimportantErr bool
}

func newLayer[T any](cache ...Cache[T]) *layer[T] {
	if len(cache) == 0 {
		return nil
	}
	this := cache[0]
	next := newLayer[T](cache[1:]...)
	return &layer[T]{
		cache: this,
		next:  next,
	}
}

func (l *layer[T]) get(ctx context.Context, key string, fn func() (T, error)) (T, error) {
	if val, err := l.cache.Get(ctx, key); err == nil {
		return val, err
	} else if !errors.Is(err, NotFindCache) {
		err = fmt.Errorf("[hybrid]cache %T get error: %w", l.cache, err)
		if l.unimportantErr {
			return *new(T), err
		} else {
			log.Printf("ERROR %s", err)
		}
	}
	var (
		val T
		err error
	)
	if l.next != nil {
		val, err = l.next.get(ctx, key, fn)
	} else {
		val, err = fn()
	}
	if err == nil || errors.Is(err, EmptyData) {
		if err2 := l.cache.Set(ctx, key, val, errors.Is(err, EmptyData)); err2 != nil {
			err = fmt.Errorf("[hybrid]cache %T set error: %w", l.cache, err2)
			if l.unimportantErr {
				return *new(T), err
			} else {
				log.Printf("ERROR %s", err)
			}
		}
	}
	return val, err
}

// Hybrid 多级缓存
type Hybrid[T any] struct {
	first *layer[T]
}

// NewHybrid 构建多级缓存
func NewHybrid[T any](cache ...Cache[T]) *Hybrid[T] {
	return &Hybrid[T]{
		first: newLayer[T](cache...),
	}
}

var (
	// NotFindCache 没找到缓存
	NotFindCache = errors.New("not find cache")
	// EmptyData 空数据
	EmptyData = errors.New("empty data")
)

// Get 获取数据，如果希望对空数据进行缓存建议返回 EmptyData
func (h *Hybrid[T]) Get(ctx context.Context, key string, fn func() (T, error)) (T, error) {
	return h.first.get(ctx, key, fn)
}

// Del 删除缓存
func (h *Hybrid[T]) Del(ctx context.Context, key string) error {
	var vErr error
	node := h.first
	for {
		if node == nil {
			break
		}
		if err := node.cache.Del(ctx, key); err != nil {
			vErr = err
		}
		node = node.next
	}
	return vErr
}
