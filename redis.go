package hybrid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type redisCache[T any] struct {
	conn redis.Conn
	ttl  time.Duration

	opt *Option
}

func (r *redisCache[T]) Set(ctx context.Context, key string, val T, isEmpty bool) error {
	ttl := r.ttl
	if isEmpty {
		if !r.opt.Empty {
			return nil
		}
		if r.opt.EmptyTtl > 0 {
			ttl = r.opt.EmptyTtl
		}
	}
	w := NewWarp(val)
	bytes, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("val serialize error:%w, [%T]%+v", err, val, val)
	}
	if _, err = r.conn.Do("SET", r.key(key), bytes, "EX", int(ttl/time.Second)); err != nil {
		return fmt.Errorf("redis set error:%w", err)
	}
	return nil
}

func (r *redisCache[T]) Get(ctx context.Context, key string) (T, error) {
	data, err := redis.Bytes(r.conn.Do("GET", r.key(key)))
	if errors.Is(err, redis.ErrNil) {
		return *new(T), NotFindCache
	}
	if err != nil {
		return *new(T), fmt.Errorf("redis get error:%w", err)
	}
	w := NewWarpEmpty[T]()
	err = json.Unmarshal(data, w)
	if err != nil {
		return *new(T), fmt.Errorf("json unmarshal error:%w", err)
	}
	return w.Value, nil
}

func (r *redisCache[T]) Del(ctx context.Context, key string) error {
	if _, err := r.conn.Do("DEL", r.key(key)); err != nil {
		return fmt.Errorf("redis del error:%w", err)
	}
	return nil
}

func (r *redisCache[T]) key(key string) string {
	if r.opt.Prefix != "" {
		return fmt.Sprintf("[%s]%s", r.opt.Prefix, key)
	}
	return key
}

// WithRedis 使用redis缓存
func WithRedis[T any](conn redis.Conn, ttl time.Duration, opts ...Options) Cache[T] {
	opt := &Option{}
	for _, optFn := range opts {
		optFn(opt)
	}
	return &redisCache[T]{
		conn: conn,
		ttl:  ttl,
		opt:  opt,
	}
}
