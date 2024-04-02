package cache

import (
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/coocood/freecache"
	"github.com/sirupsen/logrus"
)

var (
	cache *freecache.Cache
)

func InitCache() {
	cacheSize := conf.Get().Master.CacheSize * 1024 * 1024 // MB
	cache = freecache.NewCache(cacheSize)
	logrus.Infof("init cache success, size: %d MB", cacheSize/1024/1024)
}

func Get() *freecache.Cache {
	return cache
}
