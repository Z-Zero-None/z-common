package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"z-common/src/v2/core/registry"
)

var _ registry.Watcher = (*watcher)(nil)

type watcher struct {
	key         string
	ctx         context.Context
	cancel      context.CancelFunc
	watchChan   clientv3.WatchChan
	watcher     clientv3.Watcher
	kv          clientv3.KV
	serviceName string
}

func newWatcher(ctx context.Context, key, name string, client *clientv3.Client) (*watcher, error) {
	w := &watcher{
		key:         key,
		watcher:     clientv3.NewWatcher(client),
		kv:          clientv3.NewKV(client),
		serviceName: name,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.watchChan = w.watcher.Watch(w.ctx, key, clientv3.WithPrefix(), clientv3.WithRev(0), clientv3.WithKeysOnly())
	err := w.watcher.RequestProgress(w.ctx)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (obj *watcher) Next() ([]*registry.Service, error) {
	select {
	case <-obj.ctx.Done():
		return nil, obj.ctx.Err()
	case <-obj.watchChan:
		return obj.getInstance()
	}
}

func (obj *watcher) getInstance() ([]*registry.Service, error) {
	resp, err := obj.kv.Get(obj.ctx, obj.key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]*registry.Service, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		if si.Name != obj.serviceName {
			continue
		}
		items = append(items, si)
	}
	return items, nil
}

func (obj *watcher) Close() error {
	obj.cancel()
	return obj.watcher.Close()
}
