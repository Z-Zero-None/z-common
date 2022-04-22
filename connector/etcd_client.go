package connector

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"strings"
	"time"
)

const EtcdPwd = "ETCD_PWD"
const EtcdUsername = "ETCD_USERNAME"

func GetETCDCli(urls string) (*clientv3.Client, error) {
	cfg := clientv3.Config{
		Endpoints:   strings.Split(urls, ","),
		DialTimeout: 5 * time.Second,
	}

	username, password := os.Getenv(EtcdPwd), os.Getenv(EtcdUsername)
	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("clientV3.New err: %v", err)
	}
	return cli, nil
}
