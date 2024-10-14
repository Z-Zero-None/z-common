package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"z-common/src/v2/core/registry"
)

type Endpoint struct {
	Host string `json:"host"`
	Port uint64 `json:"port"`
}

type ConfigOptions struct {
	Namespace string      `json:"namespace"`
	Endpoints []*Endpoint `json:"endpoints"`
	Username  string      `json:"username"`
	Password  string      `json:"password"`
}

type options struct {
	weight  float64
	cluster string
	group   string
	scheme  string
}

type Option func(o *options)

func WithWeight(weight float64) Option {
	return func(o *options) { o.weight = weight }
}

func WithCluster(cluster string) Option {
	return func(o *options) { o.cluster = cluster }
}

func WithGroup(group string) Option {
	return func(o *options) { o.group = group }
}

func WithDefaultKind(kind string) Option {
	return func(o *options) { o.scheme = kind }
}

func serviceName(name, schema string) string {
	return fmt.Sprintf("%s.%s", name, schema)
}
func endpoint(kind, ip string, port uint64) string {
	return fmt.Sprintf("%s://%s:%d", kind, ip, port)
}
func toRegistryService(in model.Instance) *registry.Service {
	scheme := "grpc"
	if k, ok := in.Metadata["scheme"]; ok {
		scheme = k
	}
	return &registry.Service{
		ID:        in.InstanceId,
		Name:      in.ServiceName,
		Version:   in.Metadata["version"],
		Metadata:  in.Metadata,
		Endpoints: []string{endpoint(scheme, in.Ip, in.Port)},
	}
}
