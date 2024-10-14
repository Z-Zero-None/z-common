package nacos

import (
	"context"
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"net"
	"net/url"
	"strconv"
	"z-common/src/v2/core/registry"
)

var _ registry.Registry = (*Registry)(nil)

type Registry struct {
	opts    options
	clients naming_client.INamingClient
}

func New(cfg *ConfigOptions, opts ...Option) (*Registry, error) {
	var sc []constant.ServerConfig
	var commonOpts []constant.ServerOption
	op := options{
		cluster: "DEFAULT",
		group:   constant.DEFAULT_GROUP,
		weight:  100,
		scheme:  "grpc",
	}
	for _, option := range opts {
		option(&op)
	}
	for _, edp := range cfg.Endpoints {
		sc = append(sc, *constant.NewServerConfig(edp.Host, edp.Port, commonOpts...))
	}
	if len(sc) == 0 {
		return nil, errors.New("data sources")
	}
	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithUsername(cfg.Username),
		constant.WithPassword(cfg.Password),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	// create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Registry{
		clients: client,
		opts:    op,
	}, nil
}

func (r *Registry) Register(ctx context.Context, svc *registry.Service) error {
	if len(svc.Name) == 0 {
		return errors.New("name empty")
	}
	for _, edp := range svc.Endpoints {
		u, err := url.Parse(edp)
		if err != nil {
			return err
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}
		p, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		var rmd map[string]string
		if svc.Metadata == nil {
			rmd = map[string]string{
				"scheme":  u.Scheme,
				"version": svc.Version,
			}
		} else {
			rmd = make(map[string]string, len(svc.Metadata)+2)
			for k, v := range svc.Metadata {
				rmd[k] = v
			}
			rmd["scheme"] = u.Scheme
			rmd["version"] = svc.Version
		}
		_, e := r.clients.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          host,
			Port:        uint64(p),
			ServiceName: serviceName(svc.Name, u.Scheme),
			Weight:      r.opts.weight,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			Metadata:    rmd,
			ClusterName: r.opts.cluster,
			GroupName:   r.opts.group,
		})
		if e != nil {
			return fmt.Errorf("RegisterInstance err %v,%v", e, edp)
		}
	}
	return nil

}

func (r *Registry) Deregister(ctx context.Context, svc *registry.Service) error {
	for _, edp := range svc.Endpoints {
		u, err := url.Parse(edp)
		if err != nil {
			return err
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}
		p, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		if _, err = r.clients.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          host,
			Port:        uint64(p),
			ServiceName: serviceName(svc.Name, u.Scheme),
			GroupName:   r.opts.group,
			Cluster:     r.opts.cluster,
			Ephemeral:   true,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) String() string {
	return "nacos"
}
