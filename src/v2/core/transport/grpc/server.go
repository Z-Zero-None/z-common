package grpc

import (
	"google.golang.org/grpc"
	"z-common/src/v2/core/registry"
)

type Server struct {
	engine       *grpc.Server
	Addr         string
	Registry     registry.Registry
	Name         string
	AccessTokens []string
	GRPCOptions  []grpc.ServerOption
}
