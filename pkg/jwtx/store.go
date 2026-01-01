package jwtx

import (
	"context"
	"time"
)

type Storer interface {
	Set(ctx context.Context, tokenStr string, expiration time.Duration) error
	Delete(ctx context.Context, tokenStr string) error
	Check(ctx context.Context, tokenStr string) (bool, error)
	Close(ctx context.Context) error
}

type storeOptions struct {
	CacheNS string
}

type StoreOption func(*storeOptions)

func WithCacheNS(ns string) StoreOption {
	return func(o *storeOptions) {
		o.CacheNS = ns
	}
}

// Cacher 定义了缓存操作的接口，提供基本的缓存设置、获取、删除等功能
type Cacher interface {
	// Set 设置缓存值
	// ctx: 上下文对象，用于控制操作的生命周期
	// ns: 命名空间，用于区分不同的缓存区域
	// key: 缓存键名
	// value: 要缓存的值
	// expiration: 可选的过期时间，如果不提供则使用默认过期时间
	// 返回 error: 操作失败时返回错误
	Set(ctx context.Context, ns, key, value string, expiration ...time.Duration) error

	// Get 获取缓存值
	// ctx: 上下文对象，用于控制操作的生命周期
	// ns: 命名空间，用于区分不同的缓存区域
	// key: 缓存键名
	// 返回 string: 缓存的值
	// 返回 bool: 缓存是否存在
	// 返回 error: 操作失败时返回错误
	Get(ctx context.Context, ns, key string) (string, bool, error)

	// Delete 删除缓存
	// ctx: 上下文对象，用于控制操作的生命周期
	// ns: 命名空间，用于区分不同的缓存区域
	// key: 缓存键名
	// 返回 error: 操作失败时返回错误
	Delete(ctx context.Context, ns, key string) error

	// Close 关闭缓存连接
	// ctx: 上下文对象，用于控制操作的生命周期
	// 返回 error: 操作失败时返回错误
	Close(ctx context.Context) error

	// Exists 检查缓存是否存在
	// ctx: 上下文对象，用于控制操作的生命周期
	// ns: 命名空间，用于区分不同的缓存区域
	// key: 缓存键名
	// 返回 bool: 缓存是否存在
	// 返回 error: 操作失败时返回错误
	Exists(ctx context.Context, ns, key string) (bool, error)
}

func NewStoreWithCache(cache Cacher, opts ...StoreOption) Storer {
	s := &storeImpl{
		c: cache,
		opts: &storeOptions{
			CacheNS: "jwt",
		},
	}
	for _, opt := range opts {
		opt(s.opts)
	}
	return s
}

type storeImpl struct {
	c    Cacher
	opts *storeOptions
}

func (s *storeImpl) Set(ctx context.Context, tokenStr string, expiration time.Duration) error {
	return s.c.Set(ctx, s.opts.CacheNS, tokenStr, tokenStr, expiration)
}

func (s *storeImpl) Delete(ctx context.Context, tokenStr string) error {
	return s.c.Delete(ctx, s.opts.CacheNS, tokenStr)
}

func (s *storeImpl) Check(ctx context.Context, tokenStr string) (bool, error) {
	return s.c.Exists(ctx, s.opts.CacheNS, tokenStr)
}

func (s *storeImpl) Close(ctx context.Context) error {
	return s.c.Close(ctx)
}
