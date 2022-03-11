package test_method

import (
	"fmt"
	"testing"

	"z-common/config"
)

func TestNewEnvViper(t *testing.T) {
	viper, err := config.NewEnvViper()
	if err != nil {
		fmt.Println("TestNewEnvViper:", err)
	}
	fmt.Println(viper.Get("ZERO_NONE"))
}

var defaultApolloConfig = config.ApolloConfig{
	AppId:         "SampleApp",
	Cluster:       "dev",
	NamespaceName: "application",
}

func TestNewApolloCache(t *testing.T) {
	connect, err := config.NewApolloCache(&defaultApolloConfig)
	if err != nil {
		fmt.Println("TestNewApolloCache:", err)
	}
	fmt.Println(connect.Get("timeout"))
}
