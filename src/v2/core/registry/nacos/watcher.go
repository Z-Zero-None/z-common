package nacos

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"z-common/src/v2/core/registry"
)

var _ registry.Watcher = (*Watcher)(nil)

type Watcher struct {
	ctx            context.Context
	cancel         context.CancelFunc
	client         naming_client.INamingClient
	subscribeParam *vo.SubscribeParam
	quitCh         chan bool
	watchChan      chan struct{}
	serviceName    string
	groupName      string
	clusters       []string
}

func newWatcher(ctx context.Context, svcName string, reg *Registry) *Watcher {
	ctx, cancel := context.WithCancel(ctx)
	w := &Watcher{
		cancel:      cancel,
		serviceName: serviceName(svcName, reg.opts.scheme),
		clusters:    []string{reg.opts.cluster},
		groupName:   reg.opts.group,
		quitCh:      make(chan bool),
		watchChan:   make(chan struct{}, 1),
	}
	go w.watch(ctx)
	return w
}

func (w *Watcher) watch(ctx context.Context) {
	w.subscribeParam = &vo.SubscribeParam{
		ServiceName: w.serviceName,
		Clusters:    w.clusters,
		GroupName:   w.groupName,
		SubscribeCallback: func(list []model.Instance, err error) {
			select {
			case w.watchChan <- struct{}{}:
			default:
			}
		},
	}
	_ = w.client.Subscribe(w.subscribeParam)
	select {
	case w.watchChan <- struct{}{}:
	default:
	}
}

func (w *Watcher) Next() ([]*registry.Service, error) {
	select {
	case <-w.quitCh:
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
	}
	res, err := w.client.GetService(vo.GetServiceParam{
		ServiceName: w.serviceName,
		GroupName:   w.groupName,
		Clusters:    w.clusters,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*registry.Service, 0, len(res.Hosts))
	for _, in := range res.Hosts {
		items = append(items, toRegistryService(in))
	}
	return items, nil
}

func (w *Watcher) Close() (err error) {
	select {
	case <-w.quitCh:
	default:
		close(w.quitCh)
		err = w.client.Unsubscribe(w.subscribeParam)
		w.cancel()
	}
	return err
}
