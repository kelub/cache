package cache

import "github.com/Sirupsen/logrus"

type CacheAside struct {
	cache Cache
	name  string
	key   interface{}
	value interface{}
}

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

func (a *CacheAside) Write(entry *logrus.Entry, key, value interface{}) {

}
