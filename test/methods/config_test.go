package methods

import (
	"testing"
	"z-common/src/v1/base/config"
)

func TestNewEnvViper(t *testing.T) {
	viper, err := config.NewEnvViper()
	if err != nil {
		t.Errorf("TestNewEnvViper.NewEnvViper err:%v", err)
	}
	t.Log("ZERO_NONE:", viper.Get("ZERO_NONE"))
}

var defaultApolloConfig = config.ApolloConfig{
	AppId:         "SampleApp",
	Cluster:       "dev",
	NamespaceName: "application",
}

func TestNewApolloCache(t *testing.T) {
	connect, err := config.NewApolloCache(&defaultApolloConfig)
	if err != nil {
		t.Errorf("TestNewApolloCache.NewApolloCache err:%v", err)
	}
	get, err := connect.Get("timeout")
	if err != nil {
		t.Errorf("TestNewApolloCache.Get err:%v", err)
	}
	t.Log("timeout:", get)
}
