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
	}
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

// Write 设置缓存数据
func (a *CacheAside) Write(entry *logrus.Entry, key, value interface{}) (err error) {
	// write to DB
	err = a.cache.WriteDB(key, value)
	if err != nil {
		entry.Errorln("Write to DB error", err)
	}
	// delete cache
	// 最后 delete cache 可以防止并发下，写 DB 这个相对长的时间内有其它修改导致
	// 缓存不一致，想对比后 delete 缓存不一致情况 概率高很多，因为 delete cache 时间很短。
	err = a.cache.DeleteCache(key)
	if err != nil {
		entry.Errorln("Delete Cache error", err)
	}
	return nil
}

// RWThrough 读写穿策略
// 只与缓存交互，无感读写 DB
type RWThrough struct {
	// cache 缓存通用接口
	cache Cache
	// name 缓存项目名，用于日志、prom 监控
	name string
	// Write 当无缓存时，不写如缓存直接写DB
	// 可以少写一次缓存，相当于把这次写缓存分摊到下一次读，
	// 很适合要求写延迟少的情况。
	// 默认为 true
	NoWriteAllocate bool
}

// Read
// check isExist from cache
// if exist read cache
// else read from DB and write cache
func (rw *RWThrough) Read(entry *logrus.Entry, key interface{}) (value interface{}, err error) {
	//isExist, err := rw.cache.IsExist(key)
	//if err != nil {
	//	return nil, err
	//}
	value, err = rw.cache.ReadCache(key)
	if err != nil {
		return nil, err
	}
	// Hit
	if value != nil{
		CacheHit.Inc([]string{rw.name})
		return value, nil
	}
	// Miss
	CacheMiss.Inc([]string{rw.name})
	value ,err = rw.cache.ReadDB(key)
	if err != nil {
		return nil, err
	}
	err = rw.cache.WriteCache(key, value)
	if err != nil {
		entry.Errorln("write to cache error")
	}
	return value, err

}

// Write
// check isExist from cache
// if exist Update cache and update DB
// else write cache and write DB
// or write DB
func (rw *RWThrough) Write(entry *logrus.Entry, key, value interface{}) (err error) {
	isExist, err := rw.cache.IsExist(key)
	if err != nil {
		return err
	}
	if isExist {
		err = rw.cache.WriteCache(key, value)
		if err != nil {
			entry.Errorln("Write to cache error", err)
		}
		err = rw.cache.WriteDB(key, value)
		if err != nil {
			entry.Errorln("Write to DB error", err)
		}
		return nil
	}
	if !rw.NoWriteAllocate {
		err = rw.cache.WriteCache(key, value)
		if err != nil {
			entry.Errorln("Write to cache error", err)
		}
		}
		err = rw.cache.WriteDB(key, value)
		if err != nil {
			entry.Errorln("Write to DB error", err)
		}
		return nil
}

// WriteBack 写回策略
// 更新数据，只更新缓存
// 通过设置脏数据，数据再次使用才会写入DB
type WriteBack struct {
	// cache 缓存通用接口
	cache Cache
	// name 缓存项目名，用于日志、prom 监控
	name string
}

// Read
func (wb *WriteBack) Read(entry *logrus.Entry, key interface{}) (value interface{}, err error) {
	value, err = wb.cache.ReadCache(key)
	if err != nil {
		return nil, err
	}
	// Hit
	if value != nil{
		CacheHit.Inc([]string{wb.name})
		return value, nil
	}
	// Miss
	CacheMiss.Inc([]string{wb.name})
	isDirty, err := wb.cache.IsDirty(value)
	if err != nil {
		return nil, err
	}
	// isDirty
	if isDirty{
		err = wb.cache.WriteDB(key,value)
		if err != nil {
			entry.Errorln("Write to DB error", err)
		}
	}
	// read DB
	value ,err = wb.cache.ReadDB(key)
	if err != nil {
		return nil, err
	}
	// mark no dirty
	err = wb.cache.Mark(key,false)
	if err != nil {
		return nil, err
	}
	return value, err
}

// Wirte
func (rb *WriteBack) Write(entry *logrus.Entry, key, value interface{}) (err error) {
	isExist, err := rb.cache.IsExist(key)
	if err != nil {
		return err
	}
	// Miss
	if !isExist{
		CacheMiss.Inc([]string{rb.name})
		isDirty, err := rb.cache.IsDirty(value)
		if err != nil {
			return  err
		}
		// isDirty
		if isDirty {
			err = rb.cache.WriteDB(key,value)
			if err != nil {
				entry.Errorln("Write to DB error", err)
			}
		}
		// read DB
		value ,err = rb.cache.ReadDB(key)
		if err != nil {
			return  err
		}
	}
	// Hit
	CacheHit.Inc([]string{rb.name})
	err = rb.cache.WriteCache(key, value)
	if err != nil {
		entry.Errorln("Write to cache error", err)
	}
	return nil
}
