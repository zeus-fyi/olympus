package proxy_anvil

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/oleiade/lane/v2"
	"github.com/puzpuzpuz/xsync/v2"
	"github.com/rs/zerolog/log"
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
	LockedRouteToSessionCacheMap = xsync.NewMapOf[string]()
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

	if route, ok := LockedSessionToRouteCacheMap.Load(sessionID); ok {
		ttl := ts.UnixTimeStampNowSec() + int(a.LockDefaultTime.Seconds())
		LockedSessionTTL.Store(sessionID, fmt.Sprintf("%d", ttl))
		if iris_redis.IrisRedisClient.Writer != nil {
			_, _ = iris_redis.IrisRedisClient.GetAndUpdateLatestSessionCacheTTLIfExists(ctx, sessionID, a.LockDefaultTime)
		}
		return route, nil
	}

	if iris_redis.IrisRedisClient.Writer != nil {
		val, err := iris_redis.IrisRedisClient.GetAndUpdateLatestSessionCacheTTLIfExists(ctx, sessionID, a.LockDefaultTime)
		if err == nil && len(val) > 0 {
			ttl := ts.UnixTimeStampNowSec() + int(a.LockDefaultTime.Seconds())
			LockedSessionTTL.Store(sessionID, fmt.Sprintf("%d", ttl))
			return val, nil
		}
	}

	i := 0
	j := 0
	pqSize := a.PriorityQueue.Size()
	for {
		route, ttl, ok := a.PriorityQueue.Pop()
		if !ok {
			return "", ErrNoRoutesAvailable
		}
		if iris_redis.IrisRedisClient.Reader != nil && iris_redis.IrisRedisClient.Writer != nil {
			exists, err := iris_redis.IrisRedisClient.DoesSessionIDExist(ctx, sessionID)
			if err != nil {
				a.PriorityQueue.Push(route, ttl)
				log.Err(err).Msg("error checking if session exists")
				return "", err
			}
			if !exists {
				err = iris_redis.IrisRedisClient.AddSessionWithTTL(ctx, sessionID, route, a.LockDefaultTime)
				if err != nil {
					a.PriorityQueue.Push(route, ttl)
					log.Err(err).Msg("error adding session to cache")
					return "", err
				}
				LockedSessionToRouteCacheMap.Store(sessionID, route)
				LockedRouteToSessionCacheMap.Store(route, sessionID)
				ttl = ts.UnixTimeStampNowSec() + int(a.LockDefaultTime.Seconds())
				newTTL := fmt.Sprintf("%d", ttl)
				LockedSessionTTL.Store(sessionID, newTTL)
				a.PriorityQueue.Push(route, ttl)
				return route, nil
			}
		}
		a.PriorityQueue.Push(route, ttl)
		if i >= int(pqSize) {
			if j > 0 {
				return "", ErrNoRoutesAvailable
			}
			minDuration := 10 * time.Millisecond
			maxDuration := 50 * time.Millisecond
			jitter := time.Duration(j) * (time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration)
			time.Sleep(jitter)
			j++
		}
		i++
	}
}
