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
		"http://anvil.78ab2d4c-82eb-4bbc-b0fb-b702639e78c0.svc.cluster.local:8545",
	}
	ts = chronos.Chronos{}
)

func InitAnvilProxy() {
	lfuCache := datastructures.New()
	SessionLocker = AnvilProxy{
		LFU:             lfuCache,
		LockDefaultTime: time.Second * 120,
		SessionRouteMap: make(map[string]int),
		RouteLockTTL:    make(map[int]int),
	}
}

// setSessionLockOnRoute sets a session lock on a route index & updates its lock time
func (a *AnvilProxy) setSessionLockOnRoute(r *Route) (*Route, error) {
	a.routeLockTTLMutex.Lock()
	defer a.routeLockTTLMutex.Unlock()
	a.LFU.Set(r.SessionID, r.Index)
	a.SessionRouteMap[r.SessionID] = r.Index
	a.RouteLockTTL[r.Index] = ts.UnixTimeStampNow() + int(a.LockDefaultTime.Seconds())
	return r, nil
}

func (a *AnvilProxy) GetSessionLockedRoute(sessionID string) (*Route, error) {
	routeIndex := a.LFU.Get(sessionID)
	if routeIndex == nil {
		return a.GetNextAvailableRouteAndAssignToSession(sessionID)
	}
	routePath := AnvilRoutes[routeIndex.(int)]
	return &Route{
		Index:     routeIndex.(int),
		SessionID: sessionID,
		Route:     routePath,
	}, nil
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
	_, leastFreqElement := a.LFU.GetLeastFrequentValue()
	ch := make(chan datastructures.Eviction, 1)
	leastFreqSessionID := leastFreqElement.(string)
	a.routeLockTTLMutex.Lock()
	defer a.routeLockTTLMutex.Unlock()
	leastFreqRouteIndex := a.SessionRouteMap[leastFreqSessionID]
	if a.RouteLockTTL[leastFreqRouteIndex] < ts.UnixTimeStampNow() {
		// if the lock has expired, then remove it from the LFU & the session map
		delete(a.SessionRouteMap, leastFreqSessionID)
		delete(a.RouteLockTTL, leastFreqRouteIndex)
		a.LFU.Evict(1)
		return a.waitForNextAvailableRoute(ch, sessionID)
	}
	return nil, errors.New("no available routes")
}

func (a *AnvilProxy) waitForNextAvailableRoute(ch chan datastructures.Eviction, sessionID string) (*Route, error) {
	a.LFU.EvictionChannel = ch
	ev := <-ch
	log.Info().Msgf("evicted: key %v, val %v", ev.Key, ev.Value)
	return &Route{
		Index:     ev.Value.(int),
		SessionID: sessionID,
		Route:     AnvilRoutes[ev.Value.(int)],
	}, nil
}
