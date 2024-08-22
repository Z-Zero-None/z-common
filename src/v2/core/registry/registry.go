package registry

import (
	"context"
	"fmt"
)

type Service struct {
	ID        string            `json:"id"`
	Version   string            `json:"version"`
	Name      string            `json:"name"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoint"`
}

func (s *Service) String() string {
	return fmt.Sprintf("%s-%s", s.Name, s.ID)
}

type Registry interface {
	Register(ctx context.Context, svc *Service) error
	Deregister(ctx context.Context, svc *Service) error
}

type Discovery interface {
	GetService(ctx context.Context, serviceName string) ([]*Service, error)
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

type Watcher interface {
	Next() ([]*Service, error)
	Close() error
}
