package main

import "z-common/demo/memory-cache/impl"

/*
实现一个内存缓存系统
可以设置值国企
限制内存大小
确保并发安全
单元测试
*/
func main() {
	cache := impl.NewMemCache()
	cache.SetMaxMemory("300MB")
	cache.Set("zzn", 1, 0)
	print(cache.Exists("zzn"))
	print(cache.Get("zzn"))
	cache.Keys()
	cache.Del("zzn")
	cache.Flush()
}
