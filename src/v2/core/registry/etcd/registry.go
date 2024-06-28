package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
	"z-common/src/v2/core/registry"
)

var (
	_ registry.Registry  = nil
	_ registry.Discovery = nil
)

type Option func(o *options)

type options struct {
	ctx       context.Context
	namespace string
	ttl       time.Duration
	maxRetry  int
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

type Registry struct {
	opts   *options
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
	ctxMap map[*registry.Service]context.CancelFunc
}
