package proxy_anvil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oleiade/lane/v2"
	"github.com/puzpuzpuz/xsync/v2"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

var (
	Routes = []string{
		"http://anvil.191aada9-055d-4dba-a906-7dfbc4e632c6.svc.cluster.local:8545",
		"http://anvil.427c5536-4fc0-4257-90b5-1789d290058c.svc.cluster.local:8545",
		"http://anvil.5cf3a2c0-1d65-48cb-8b85-dc777ad956a0.svc.cluster.local:8545",
		"http://anvil.78ab2d4c-82eb-4bbc-b0fb-b702639e78c0.svc.cluster.local:8545",
		"http://anvil.a49ca82d-ff96-4c4f-8653-001d56cab5e5.svc.cluster.local:8545",
		"http://anvil.be58f278-1fbe-4bc8-8db5-03d8901cc060.svc.cluster.local:8545",
		"http://anvil.e56def19-190f-4b45-9fdb-8468ddbe0eb5.svc.cluster.local:8545",
	}
	ts                   = chronos.Chronos{}
	SessionLocker        = AnvilProxy{}
	ErrNoRoutesAvailable = errors.New("no routes available")
)

type AnvilProxy struct {
	LockDefaultTime time.Duration
	PriorityQueue   *lane.PriorityQueue[string, int]
}

func InitAnvilProxy() {
	SessionLocker.PriorityQueue = lane.NewMinPriorityQueue[string, int]()
	SessionLocker.LockDefaultTime = 30 * time.Second
	for _, route := range Routes {
		SessionLocker.PriorityQueue.Push(route, -1)
	}
}

var (
	LockedSessionTTL             = xsync.NewMapOf[string]()
	LockedSessionToRouteCacheMap = xsync.NewMapOf[string]()
)

func (a *AnvilProxy) RemoveSessionLockedRoute(ctx context.Context, sessionID string) {
	if _, ok := LockedSessionToRouteCacheMap.Load(sessionID); ok {
		LockedSessionTTL.Store(sessionID, fmt.Sprintf("%d", 0))
	}
	if iris_redis.IrisRedisClient.Writer != nil {
		_, _ = iris_redis.IrisRedisClient.DeleteSessionCacheIfExists(ctx, sessionID)
		_, _ = iris_redis.IrisRedisClient.DeleteSessionRouteCacheIfExists(ctx, sessionID)
	}
	return
}

func (a *AnvilProxy) GetSessionLockedRoute(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "Zeus-Test" {
		return "http://anvil.eeb335ad-78da-458f-9cfb-9928514d65d0.svc.cluster.local:8545", nil
	}
	val, err := iris_redis.IrisRedisClient.GetAndUpdateLatestSessionCacheTTLIfExists(ctx, sessionID, a.LockDefaultTime)
	if err == nil && len(val) > 0 {
		return val, nil
	}
	return "", err
}
