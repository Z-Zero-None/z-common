package test_config

import (
	"fmt"
	"testing"

	"z-common/config"
)

func TestNewEnvViper(t *testing.T) {
	viper, err := config.NewEnvViper()
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(viper.Get("ZERO_NONE"))
}