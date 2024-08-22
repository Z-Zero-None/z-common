package etcd

import (
	"encoding/json"
	"z-common/src/v2/core/registry"
)

func marshal(si *registry.Service) (string, error) {
	data, err := json.Marshal(si)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(data []byte) (si *registry.Service, err error) {
	err = json.Unmarshal(data, &si)
	return
}
