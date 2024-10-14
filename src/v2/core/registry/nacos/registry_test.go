package nacos

import (
	"context"
	"github.com/google/uuid"
	"log"
	"testing"
	"z-common/src/v2/core/registry"
)

func TestNew(t *testing.T) {
	cfg := &ConfigOptions{
		Namespace: "public",
		Endpoints: []*Endpoint{
			{
				Host: "127.0.0.1",
				Port: 8848,
			},
		},
		Username: "nacos",
		Password: "123456789",
	}
	r, err := New(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	userSvc := &registry.Service{
		ID:      uuid.New().String(),
		Version: "v1",
		Name:    "users",
		Metadata: map[string]string{
			"name": "one two",
		},
		Endpoints: []string{"grpc://10.255.255.254:11170"},
	}
	ctx := context.Background()
	log.Printf("regiter err %v", r.Register(ctx, userSvc))
	if list, err := r.GetService(ctx, "users.grpc"); err == nil {
		for _, item := range list {
			log.Printf("item %#v", item)
		}

	} else {
		log.Fatalln(err)
	}
}
