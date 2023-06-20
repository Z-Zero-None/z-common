package connector

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var RocketMQProducer rocketmq.Producer

var defaultEndPoints = []string{"127.0.0.1:9876"}
var defaultTopic = "test"

func InitProducer(endPoints ...string) error {
	if len(endPoints) < 1 {
		endPoints = defaultEndPoints
	}
	addr, err := primitive.NewNamesrvAddr(endPoints...)
	if err != nil {
		return err
	}
	// 注册生产者
	if RocketMQProducer, err = rocketmq.NewProducer(
		producer.WithNameServer(addr),
		producer.WithRetry(3),
	); err != nil {
		return err
	}
	return nil
}
