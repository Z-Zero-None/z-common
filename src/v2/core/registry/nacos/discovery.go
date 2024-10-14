package nacos

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"z-common/src/v2/core/registry"
)

var _ registry.Discovery = (*Registry)(nil)

func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.Service, error) {
	res, err := r.clients.SelectInstances(vo.SelectInstancesParam{
		ServiceName: name,
		GroupName:   r.opts.group,
		HealthyOnly: true,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*registry.Service, 0, len(res))
	for _, in := range res {
		items = append(items, toRegistryService(in))
	}
	return items, nil
}

func (r *Registry) List(ctx context.Context) ([]*registry.Service, error) {
	var items []*registry.Service
	page := 1
	for {
		svcNames, err := r.clients.GetAllServicesInfo(vo.GetAllServiceInfoParam{
			PageNo:   uint32(page),
			PageSize: 10,
		})
		if err != nil {
			return nil, err
		}
		if len(svcNames.Doms) == 0 {
			break
		}
		page++
		for _, name := range svcNames.Doms {
			res, err := r.clients.SelectAllInstances(vo.SelectAllInstancesParam{
				ServiceName: name,
				GroupName:   r.opts.group,
			})
			if err != nil {
				return nil, err
			}
			for _, in := range res {
				items = append(items, toRegistryService(in))
			}
		}
	}

	return items, nil
}

func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	return newWatcher(ctx, name, r), nil
}
