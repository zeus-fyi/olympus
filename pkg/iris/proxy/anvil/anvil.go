package proxy_anvil

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type AnvilProxy struct {
}

var (
	AnvilRoutes            = make(map[string]string)
	AnvilLocalSessionCache = cache.New(1*time.Second, 2*time.Second)
)

func AddToAnvilLocalSessionCache(key string, value interface{}) {
	AnvilLocalSessionCache.Set(key, value, cache.DefaultExpiration)
}

func GetFromAnvilLocalSessionCache(key string) (interface{}, bool) {
	return AnvilLocalSessionCache.Get(key)
}

func DeleteFromAnvilLocalSessionCache(key string) {
	AnvilLocalSessionCache.Delete(key)
}

func (a *AnvilProxy) SetSessionLock() {

}
