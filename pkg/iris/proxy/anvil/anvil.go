package proxy_anvil

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/oleiade/lane/v2"
	"github.com/puzpuzpuz/xsync/v2"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

var (
	LockedSessionTTL             = xsync.NewMapOf[string]()
	LockedSessionToRouteCacheMap = xsync.NewMapOf[string]()
	LockedRouteToSessionCacheMap = xsync.NewMapOf[string]()
	Routes                       = []string{
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
	for _, route := range Routes {
		SessionLocker.PriorityQueue.Push(route, -1)
	}
}

func (a *AnvilProxy) RemoveSessionLockedRoute(sessionID string) {
	if _, ok := LockedSessionToRouteCacheMap.Load(sessionID); ok {
		LockedSessionTTL.Store(sessionID, fmt.Sprintf("%d", 0))
		return
	}
	return
}

func (a *AnvilProxy) GetSessionLockedRoute(ctx context.Context, sessionID string) (string, error) {
	if route, ok := LockedSessionToRouteCacheMap.Load(sessionID); ok {
		ttl := ts.UnixTimeStampNowSec() + int(a.LockDefaultTime.Seconds())
		LockedSessionTTL.Store(sessionID, fmt.Sprintf("%d", ttl))
		return route, nil
	}
	for {
		route, ttl, ok := a.PriorityQueue.Pop()
		if !ok {
			return "", ErrNoRoutesAvailable
		}
		oldSession, exists := LockedRouteToSessionCacheMap.Load(route)
		if exists {
			mapTTL, mapTTLExists := LockedSessionTTL.Load(oldSession)
			if mapTTLExists {
				ttlLatest, err := strconv.Atoi(mapTTL)
				if err != nil {
					log.Err(err).Msg("error converting ttl to int")
					return "", err
				}
				ttl = ttlLatest
			}
		}
		if ttl < ts.UnixTimeStampNowSec() {
			if exists {
				LockedSessionToRouteCacheMap.Delete(oldSession)
				LockedRouteToSessionCacheMap.Delete(route)
				LockedSessionTTL.Delete(oldSession)
			}
			LockedSessionToRouteCacheMap.Store(sessionID, route)
			LockedRouteToSessionCacheMap.Store(route, sessionID)
			ttl = ts.UnixTimeStampNowSec() + int(a.LockDefaultTime.Seconds())
			newTTL := fmt.Sprintf("%d", ttl)
			LockedSessionTTL.Store(sessionID, newTTL)
			a.PriorityQueue.Push(route, ttl)
			return route, nil
		}
		a.PriorityQueue.Push(route, ttl)
	}
}
