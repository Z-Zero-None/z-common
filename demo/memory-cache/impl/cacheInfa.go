package impl

import "time"

type CacheInfa interface {
	//SetMaxMemory 设置内存大小
	SetMaxMemory(size string) bool // 设置内存大小
	//Set 写入缓存
	Set(key string, val interface{}, expire time.Duration) bool
	//Get 依据key获取数据
	Get(key string) (interface{}, bool)
	//Del 删除key
	Del(key string) bool
	//Exists 是否存在key
	Exists(key string) bool
	//Flush 清空
	Flush() bool
	//Keys 获取key的总数
	Keys() int64
}
