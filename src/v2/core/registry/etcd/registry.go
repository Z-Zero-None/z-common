package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math/rand"
	"time"
	"z-common/src/v2/core/registry"
)

var (
	_ registry.Registry  = nil
	_ registry.Discovery = nil
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

type Registry struct {
	opts   *options
	client *clientv3.Client //etcd客户端
	kv     clientv3.KV      //keys-value操作
	lease  clientv3.Lease
	ctxMap map[*registry.Service]context.CancelFunc //停止心跳处理
}

func (r *Registry) key(name, id string) string {
	return fmt.Sprintf("%s/%s/%s", r.opts.namespace, name, id)
}

func (r *Registry) serviceKey(name string) string {
	return fmt.Sprintf("%s/%s", r.opts.namespace, name)
}

func (r *Registry) Register(ctx context.Context, svc *registry.Service) error {
	key := r.key(svc.Name, svc.ID)
	value, err := marshal(svc)
	if err != nil {
		return err
	}
	if r.lease != nil {
		r.lease.Close()
	}
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.register(ctx, key, value)
	if err != nil {
		return err
	}

	hctx, cancel := context.WithCancel(r.opts.ctx)
	r.ctxMap[svc] = cancel
	go r.heartbeat(hctx, leaseID, key, value)
	return nil
}

func (r *Registry) heartbeat(ctx context.Context, leaseID clientv3.LeaseID, key, val string) {
	curLeaseID := leaseID
	kac, err := r.client.KeepAlive(ctx, leaseID)
	if err != nil {
		curLeaseID = 0
	}
	rand.New(rand.NewSource(time.Now().Unix()))
	for {
		if curLeaseID == 0 {
			// try to registerWithKV
			var retreat []int
			for retryCnt := 0; retryCnt < r.opts.maxRetry; retryCnt++ {
				if ctx.Err() != nil {
					return
				}
				// prevent infinite blocking
				idChan := make(chan clientv3.LeaseID, 1)
				errChan := make(chan error, 1)
				cancelCtx, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()
					id, registerErr := r.register(cancelCtx, key, val)
					if registerErr != nil {
						errChan <- registerErr
					} else {
						idChan <- id
					}
				}()

				select {
				case <-time.After(3 * time.Second):
					cancel()
					continue
				case <-errChan:
					continue
				case curLeaseID = <-idChan:
				}

				kac, err = r.client.KeepAlive(ctx, curLeaseID)
				if err == nil {
					break
				}
				retreat = append(retreat, 1<<retryCnt)
				time.Sleep(time.Duration(retreat[rand.Intn(len(retreat))]) * time.Second)
			}
			if _, ok := <-kac; !ok {
				// retry failed
				return
			}
		}

		select {
		case _, ok := <-kac:
			if !ok {
				if ctx.Err() != nil {
					// channel closed due to context cancel
					return
				}
				// need to retry registration
				curLeaseID = 0
				continue
			}
		case <-r.opts.ctx.Done():
			return
		}
	}
}

func (r *Registry) register(ctx context.Context, key, val string) (clientv3.LeaseID, error) {
	grant, err := r.lease.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return 0, err
	}
	_, err = r.client.Put(ctx, key, val, clientv3.WithLease(grant.ID))
	if err != nil {
		return 0, err
	}
	return grant.ID, nil
}

func (r *Registry) Deregister(ctx context.Context, svc *registry.Service) error {
	defer func() {
		if r.lease != nil {
			r.lease.Close()
		}
	}()
	// cancel heartbeat
	if cancel, ok := r.ctxMap[svc]; ok {
		cancel()
		delete(r.ctxMap, svc)
	}
	key := r.key(svc.Name, svc.ID)
	_, err := r.client.Delete(ctx, key)
	return err
}

func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.Service, error) {
	key := r.serviceKey(name)
	resp, err := r.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]*registry.Service, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		if si.Name != name {
			continue
		}
		items = append(items, si)
	}
	return items, nil
}

func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	//key := r.serviceKey(name)
	//return newWatcher(ctx, key, name, r.client)
	return nil, nil
}
