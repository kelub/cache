package cache

import (
	"github.com/Sirupsen/logrus"
)

// CacheAside 边路策略
type CacheAside struct {
	// cache 缓存通用接口
	cache Cache
	// name 缓存项目名，用于日志、prom 监控
	name string
}

// Read 获取缓存数据
func (a *CacheAside) Read(entry *logrus.Entry, key interface{}) (value interface{}, err error) {
	//read from cache
	value, err = a.cache.ReadCache(key)
	if err != nil {
		return nil, err
	}
	//hit return date
	if value != nil {
		CacheHit.Inc([]string{a.name})
		return value, nil
	} else {
		// miss
		// read from DB
		CacheMiss.Inc([]string{a.name})
		value, err = a.cache.ReadDB(key)
		if err != nil {
			return nil, err
		}
		// write to cache
		err = a.cache.WriteCache(key, value)
		if err != nil {
			entry.Errorln("write to cache error")
		}
		return value, err
	}
}

// Write 设置缓存数据
func (a *CacheAside) Write(entry *logrus.Entry, key, value interface{}) (err error) {
	//write to DB
	err = a.cache.WriteDB(key, value)
	if err != nil {
		entry.Errorln("Write to DB error", err)
	}
	//delete cache
	err = a.cache.DeleteCache(key)
	if err != nil {
		entry.Errorln("Delete Cache error", err)
	}
	return nil
}
