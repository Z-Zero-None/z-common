package methods

import (
	"testing"
	config2 "z-common/src/base/config"
)

func TestNewEnvViper(t *testing.T) {
	viper, err := config2.NewEnvViper()
	if err != nil {
		t.Errorf("TestNewEnvViper.NewEnvViper err:%v", err)
	}
	t.Log("ZERO_NONE:", viper.Get("ZERO_NONE"))
}

var defaultApolloConfig = config2.ApolloConfig{
	AppId:         "SampleApp",
	Cluster:       "dev",
	NamespaceName: "application",
}

func TestNewApolloCache(t *testing.T) {
	connect, err := config2.NewApolloCache(&defaultApolloConfig)
	if err != nil {
		t.Errorf("TestNewApolloCache.NewApolloCache err:%v", err)
	}
	get, err := connect.Get("timeout")
	if err != nil {
		t.Errorf("TestNewApolloCache.Get err:%v", err)
	}
	t.Log("timeout:", get)
}
