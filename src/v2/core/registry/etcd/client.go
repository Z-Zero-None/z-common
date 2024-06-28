package etcd

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func NewETCD(cfg *viper.Viper) (*clientv3.Client, error) {
	etcdCfg := clientv3.Config{
		Endpoints:   cfg.GetStringSlice("etcd.endpoints"),
		DialTimeout: 5 * time.Second,
		Username:    cfg.GetString("etcd.user"),
		Password:    cfg.GetString("etcd.password"),
		//TLS: nil,
	}

	return clientv3.New(etcdCfg)
}
