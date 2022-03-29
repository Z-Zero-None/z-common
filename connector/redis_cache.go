package connector

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisConfig struct {
	Host        string `json:"host"`
	Password    string `json:"password"`
	MaxIdle     int    `json:"maxIdle"`
	MaxActive   int    `json:"maxActive"`
	IdleTimeout int    `json:"idleTimeout"`
	DB          int    `json:"db"`
	Wait        bool   `json:"wait"`
}

var defaultRedisConfig = RedisConfig{
	Host:        "localhost:6379",
	Password:    "",
	MaxIdle:     8,
	MaxActive:   8,
	IdleTimeout: 240,
	DB:          0,
	Wait:        false,
}

func NewDefaultRedisConfig() *RedisConfig {
	return &defaultRedisConfig
}

func redisPing(c redis.Conn, t time.Time) error {
	_, err := c.Do("PING")
	return err
}

func GetRedisCachePool(rc *RedisConfig) (rp *redis.Pool, err error) {
	if rc == nil {
		return nil, fmt.Errorf("RedisConfig is nil")
	}
	if rc.Host == "" {
		return nil, fmt.Errorf("lack of RedisConfig.Host")
	}
	//设置默认值
	if !(rc.MaxActive > 0) {
		rc.MaxActive = 8
	}
	if !(rc.MaxIdle > 0) {
		rc.MaxIdle = 8
	}
	if !(rc.IdleTimeout > 0) {
		rc.IdleTimeout = 240
	}
	var dialOptions []redis.DialOption
	if len(rc.Password) > 0 {
		dialOptions = append(dialOptions, redis.DialPassword(rc.Password))
	}
	if rc.DB > 0 {
		dialOptions = append(dialOptions, redis.DialDatabase(rc.DB))
	}
	rp = &redis.Pool{
		MaxIdle:     rc.MaxIdle,
		MaxActive:   rc.MaxActive,
		IdleTimeout: time.Duration(rc.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rc.Host, dialOptions...)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: redisPing,
	}
	return rp, nil
}
