package impl

import (
	"sync"
	"time"

	"github.com/apolloconfig/agollo/v4/component/log"
)

var defaultSize = 1024

type MemCache struct {
	// 最大内存
	maxMemorySize int64
	//内存字符串表示
	maxMemorySizeStr string
	// 当前使用内存大小
	currMemorySize int64
	//存储值
	values map[string]*memCacheValue
	//声明锁变量
	lock sync.RWMutex
	//清楚过期数据周期
	clearTime time.Duration
}

type memCacheValue struct {
	val  interface{}
	exp  time.Time
	size int64
}

func NewMemCache() CacheInfa {
	m := &MemCache{
		values:    make(map[string]*memCacheValue, defaultSize),
		clearTime: time.Minute * 5,
	}
	go m.cleanCron()

	return m
}

func (m *MemCache) SetMaxMemory(size string) bool {
	m.maxMemorySize, m.maxMemorySizeStr = parseSize(size)
	println("SetMaxMemory:", m.maxMemorySize, "	", m.maxMemorySizeStr)
	return false
}

// 写入缓存
func (m *MemCache) Set(key string, val interface{}, expire time.Duration) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	v := &memCacheValue{
		val:  val,
		exp:  time.Now().Add(expire),
		size: getValSize(val),
	}
	m.del(key)
	m.add(key, v)
	if m.currMemorySize > m.maxMemorySize {
		m.del(key)
		//可以替换为清理过期数据
		//panic(fmt.Sprintf("max memory size %s", m.maxMemorySize))
		log.Errorf("max memory size %s\n", m.maxMemorySize)
		return false
	}
	return true
}
func (m *MemCache) get(key string) (*memCacheValue, bool) {
	value, ok := m.values[key]
	return value, ok
}

func (m *MemCache) del(key string) {
	tmp, ok := m.get(key)
	if ok && tmp != nil {
		m.currMemorySize -= tmp.size
		delete(m.values, key)
	}
}

func (m *MemCache) add(key string, val *memCacheValue) {
	m.values[key] = val
	m.currMemorySize += val.size
}

// 依据key获取数据
func (m *MemCache) Get(key string) (interface{}, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	v, ok := m.get(key)
	if ok {
		if v.exp.Before(time.Now()) {
			m.del(key)
			return nil, false
		}
		return v.val, true
	}
	return nil, false
}

// 删除key
func (m *MemCache) Del(key string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.del(key)
	return true
}

// 是否存在key
func (m *MemCache) Exists(key string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.get(key)
	return ok
}

// 清空
func (m *MemCache) Flush() bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	//例用go的回收机制
	m.values = make(map[string]*memCacheValue, 0)
	m.currMemorySize = 0
	return true
}

// 获取key的总数
func (m *MemCache) Keys() int64 {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return int64(len(m.values))
}

func (m *MemCache) cleanCron() {
	ticker := time.NewTicker(m.clearTime)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for key, val := range m.values {
				if val.exp.IsZero() && time.Now().After(val.exp) {
					m.lock.Lock()
					m.del(key)
					m.lock.Unlock()
				}
			}
		}
	}
}
