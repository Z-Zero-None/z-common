package registry

import (
	"context"
	"fmt"
)

type Service struct {
	ID        string            `json:"id"`
	Version   string            `json:"version"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoint"`
}

func (s *Service) String() string {
	return fmt.Sprintf("%s-%s", s.Name, s.ID)
}

type Watcher interface {
	Next(ctx context.Context) ([]*Service, error)
	Cancel() error
}

type Registry interface {
	Register(ctx context.Context, svc *Service) error
	Deregister(ctx context.Context, svc *Service) error
}

type Discovery interface {
	GetService(ctx context.Context, serviceName string) ([]*Service, error)
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}
