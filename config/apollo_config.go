package config

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache"
	config "github.com/apolloconfig/agollo/v4/env/config"
	"os"
)

const ApolloSecretKey = "APOLLO_SECRET_KEY"
const ApolloIpKey = "APOLLO_IP_KEY"

type ApolloConfig struct {
	AppId         string
	Cluster       string
	NamespaceName string
}

func NewApolloCache(con *ApolloConfig) (agcache.CacheInterface, error) {
	c := &config.AppConfig{
		AppID:   con.AppId,
		Cluster: con.Cluster,
		//"http://127.0.0.1:8080"
		IP:             os.Getenv(ApolloIpKey),
		NamespaceName:  con.NamespaceName,
		IsBackupConfig: true,
		Secret:         os.Getenv(ApolloSecretKey),
	}
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		return nil, err
	}
	//Use your apollo key to test
	cache := client.GetConfigCache(c.NamespaceName)
	return cache, nil
}
