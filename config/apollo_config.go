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

func NewApolloCache(con *ApolloConfig) agcache.CacheInterface {
	c := &config.AppConfig{
		AppID:          con.AppId,
		Cluster:        con.Cluster,
		IP:             os.Getenv(ApolloIpKey),
		NamespaceName:  con.NamespaceName,
		IsBackupConfig: true,
		Secret:         os.Getenv(ApolloSecretKey),
	}
	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	//Use your apollo key to test
	cache := client.GetConfigCache(c.NamespaceName)
	return cache
}
