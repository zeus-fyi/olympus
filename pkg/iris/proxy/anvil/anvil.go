package proxy_anvil

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
)

type AnvilProxy struct {
}

var (
	AnvilRoutes = map[string]string{
		"hardhat":  "https://hardhat.zeus.fyi",
		"anvilOne": "http://anvil.191aada9-055d-4dba-a906-7dfbc4e632c6.svc.cluster.local:8545",
		"anvilTwo": "http://anvil.78ab2d4c-82eb-4bbc-b0fb-b702639e78c0.svc.cluster.local:8545",
	}
	AnvilLocalSessionCache = cache.New(120*time.Second, 240*time.Second)
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

func (a *AnvilProxy) SetSessionLock(ctx context.Context, lockID string) {

}
