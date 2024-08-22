package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"z-common/src/v2/core/registry"
)

var _ registry.Discovery = nil

func (r *Registry) watchKey(name string) string {
	return fmt.Sprintf("%s/%s", r.opts.namespace, name)
}
func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.Service, error) {
	key := r.watchKey(name)
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
	key := r.watchKey(name)
	return newWatcher(ctx, key, name, r.client)
}
