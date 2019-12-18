package cache

import (
	"github.com/Sirupsen/logrus"
	prom "kelub/promserver/stats"
)

var (
	CacheHit  = prom.NewPromVec("cache").Counter("cache_hit", []string{"name"})
	CacheMiss = prom.NewPromVec("cache").Counter("cache_miss", []string{"name"})
)

// Cache 缓存接口
type Cache interface {
	WriteCache(key, value interface{}) error
	ReadCache(key interface{}) (value interface{}, err error)
	DeleteCache(key interface{}) error

	WriteDB(key, value interface{}) error
	ReadDB(key interface{}) (value interface{}, err error)

	// IsExist olny RWThrough/WriteBack Strategy
	IsExist(key interface{}) (isExist bool, err error)

	// IsDirty olny WriteBack Strategy
	IsDirty(key interface{}) (isDirty bool, err error)

	// Mark olny WriteBack Strategy
	// when isDirty true,mark is dirty,false mark is no dirty
	Mark(key interface{},isDirty bool) (err error)
}

type CacheStrategy interface {
	Read(entry *logrus.Entry, key interface{})
	Write(entry *logrus.Entry, key, value interface{})

	//Delete()
}
