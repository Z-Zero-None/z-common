package test_method

import (
	"github.com/gomodule/redigo/redis"
	"testing"
	"z-common/connector"
)

func TestGetMySQLEngine(t *testing.T) {
	config := connector.NewDefaultMysqlConfig()
	_, err := config.GetMySQLEngine()
	if err != nil {
		t.Errorf("TestGetMySQLEngine.GetMySQLEngine err:%v", err)
	}
	t.Log("DB success")
}

func BenchmarkGetDsnByString(b *testing.B) {
	config := connector.NewDefaultMysqlConfig()
	for i := 0; i < b.N; i++ {
		config.GetDsnByString()
	}
}

func BenchmarkGetDsnByBuffer(b *testing.B) {
	config := connector.NewDefaultMysqlConfig()
	for i := 0; i < b.N; i++ {
		config.GetDsnByString()
	}
}

//相关文档 https://pkg.go.dev/github.com/gomodule/redigo/redis#SlowLog
func TestGetRedisCachePool(t *testing.T) {
	config := connector.NewDefaultRedisConfig()
	pool, err := connector.GetRedisCachePool(config)
	if err != nil {
		t.Errorf("TestGetRedisCachePool.GetRedisCachePool err:%v", err)
	}
	cache := pool.Get()
	_, err = cache.Do("set", "zzn", 1)
	if err != nil {
		t.Errorf("TestGetRedisCachePool.Do.set err:%v", err)
	}
	reply, err := redis.String(cache.Do("get", "zzn"))
	if err != nil {
		t.Errorf("TestGetRedisCachePool.Do.get err:%v", err)
	}
	t.Log("key:zzn,value:", reply)
}

func TestGetETCDCli(t *testing.T) {
	_, err := connector.GetETCDCli("127.0.0.1:2379")
	if err != nil {
		t.Errorf("TestGetETCDCli.GetETCDCli err:%v", err)
	}
	t.Log("DB success")
}

func TestNewJaegerTrace(t *testing.T) {
	config := connector.NewDefaultJaegerTraceConfig()
	_, _, err := connector.NewJaegerTrace(config)
	if err != nil {
		t.Errorf("TestNewJaegerTrace.NewJaegerTrace err:%v", err)
	}
	t.Log("jaeger connect success")
}
