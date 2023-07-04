package proxy_anvil

import (
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/datastructures"
)

var SessionLocker = AnvilProxy{}

type AnvilProxy struct {
	LFU               *datastructures.LFU
	SessionRouteMap   map[string]int
	RouteLockTTL      map[int]int
	routeLockTTLMutex sync.RWMutex // Add a mutex to protect the map
	LockDefaultTime   time.Duration
}

type Route struct {
	Index     int
	SessionID string
	Route     string
}

var (
	AnvilRoutes = []string{
		"http://anvil.191aada9-055d-4dba-a906-7dfbc4e632c6.svc.cluster.local:8545",
		"http://anvil.427c5536-4fc0-4257-90b5-1789d290058c.svc.cluster.local:8545",
		"http://anvil.5cf3a2c0-1d65-48cb-8b85-dc777ad956a0.svc.cluster.local:8545",
		"http://anvil.78ab2d4c-82eb-4bbc-b0fb-b702639e78c0.svc.cluster.local:8545",
		"http://anvil.a49ca82d-ff96-4c4f-8653-001d56cab5e5.svc.cluster.local:8545",
		"http://anvil.be58f278-1fbe-4bc8-8db5-03d8901cc060.svc.cluster.local:8545",
		"http://anvil.e56def19-190f-4b45-9fdb-8468ddbe0eb5.svc.cluster.local:8545",
	}
	ts = chronos.Chronos{}
)

func InitAnvilProxy() {
	lfuCache := datastructures.New()
	SessionLocker = AnvilProxy{
		LFU:               lfuCache,
		LockDefaultTime:   time.Second * 6,
		SessionRouteMap:   make(map[string]int),
		RouteLockTTL:      make(map[int]int),
		routeLockTTLMutex: sync.RWMutex{},
	}
}

// setSessionLockOnRoute sets a session lock on a route index & updates its lock time
func (a *AnvilProxy) setSessionLockOnRoute(r *Route) (*Route, error) {
	a.routeLockTTLMutex.Lock()
	defer a.routeLockTTLMutex.Unlock()
	a.LFU.Set(r.SessionID, r.Index)
	a.SessionRouteMap[r.SessionID] = r.Index
	a.RouteLockTTL[r.Index] = ts.UnixTimeStampNowSec() + int(a.LockDefaultTime.Seconds())
	return r, nil
}

func (a *AnvilProxy) RemoveSessionLockedRoute(sessionID string) {
	log.Info().Msgf("Removing session lock for sessionID: %s", sessionID)
	leastFreqRouteIndex := a.SessionRouteMap[sessionID]
	a.RouteLockTTL[leastFreqRouteIndex] = 0
}

func (a *AnvilProxy) GetSessionLockedRoute(sessionID string) (*Route, error) {
	routeIndex := a.LFU.Get(sessionID)
	if routeIndex == nil {
		r, err := a.GetNextAvailableRouteAndAssignToSession(sessionID)
		if err != nil {
			log.Err(err).Msg("error getting next available route")
			return r, err
		}
	}
	routePath := AnvilRoutes[routeIndex.(int)]
	r := &Route{
		Index:     routeIndex.(int),
		SessionID: sessionID,
		Route:     routePath,
	}
	a.SessionRouteMap[r.SessionID] = r.Index
	a.RouteLockTTL[r.Index] = ts.UnixTimeStampNowSec() + int(a.LockDefaultTime.Seconds())
	return r, nil
}

func (a *AnvilProxy) GetNextAvailableRouteAndAssignToSession(sessionID string) (*Route, error) {
	if a.LFU.Len() < len(AnvilRoutes) {
		r := &Route{
			Index:     a.LFU.Len(),
			SessionID: sessionID,
			Route:     AnvilRoutes[a.LFU.Len()],
		}
		return a.setSessionLockOnRoute(r)
	}
	leastFreqElement, _ := a.LFU.GetLeastFrequentValue()
	leastFreqSessionID := leastFreqElement.(string)
	a.routeLockTTLMutex.Lock()
	defer a.routeLockTTLMutex.Unlock()
	leastFreqRouteIndex := a.SessionRouteMap[leastFreqSessionID]
	if a.RouteLockTTL[leastFreqRouteIndex] < ts.UnixTimeStampNowSec() {
		// if the lock has expired, then remove it from the LFU & the session map
		delete(a.SessionRouteMap, leastFreqSessionID)
		delete(a.RouteLockTTL, leastFreqRouteIndex)
		return a.waitForNextAvailableRoute(sessionID)
	}
	return nil, errors.New("no available routes")
}

func (a *AnvilProxy) waitForNextAvailableRoute(sessionID string) (*Route, error) {
	ch := make(chan datastructures.Eviction, 1)
	a.LFU.EvictionChannel = ch
	a.LFU.Evict(1)
	ev := <-ch
	log.Info().Msgf("evicted: key %v, val %v", ev.Key, ev.Value)
	r := &Route{
		Index:     ev.Value.(int),
		SessionID: sessionID,
		Route:     AnvilRoutes[ev.Value.(int)],
	}
	a.LFU.Set(r.SessionID, r.Index)
	a.SessionRouteMap[r.SessionID] = r.Index
	a.RouteLockTTL[r.Index] = ts.UnixTimeStampNowSec() + int(a.LockDefaultTime.Seconds())
	return r, nil
}
