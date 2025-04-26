package cache

import (
	"context"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/coocood/freecache"
)

var (
	cache *freecache.Cache
)

func InitCache(cfg conf.Config) {
	c := context.Background()
	cacheSize := cfg.Master.CacheSize * 1024 * 1024 // MB
	cache = freecache.NewCache(cacheSize)
	logger.Logger(c).Infof("init cache success, size: %d MB", cacheSize/1024/1024)
}

func Get() *freecache.Cache {
	return cache
}
