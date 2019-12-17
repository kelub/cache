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
}

type CacheStrategy interface {
	Read(entry *logrus.Entry, key interface{})
	Write(entry *logrus.Entry, key, value interface{})

	//Delete()
}
