package etcd

import (
	"context"
	"time"
)

type Option func(o *options)

type options struct {
	ctx       context.Context //上下文
	namespace string
	ttl       time.Duration
	maxRetry  int //重试次数
}

func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

func Namespace(ns string) Option {
	return func(o *options) { o.namespace = ns }
}

func RegisterTTL(ttl time.Duration) Option {
	return func(o *options) { o.ttl = ttl }
}

func MaxRetry(num int) Option {
	return func(o *options) { o.maxRetry = num }
}
