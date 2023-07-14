package iris_round_robin

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	RoundRobinCache      = cache.New(5*time.Minute, 10*time.Minute)
	RoundRobinRouteTable = cache.New(cache.NoExpiration, cache.NoExpiration)
)

func GetNextRoute(orgID int, rgName string) string {
	tag := fmt.Sprintf("%d-%s", orgID, rgName)
	routeTable, ok := RoundRobinRouteTable.Get(tag)
	if !ok {
		panic("no route table found")
	}
	nr, ok := RoundRobinCache.Get(rgName)
	if !ok {
		panic("no route group found")
	}
	nr = (nr.(int) + 1) % len(routeTable.([]string))
	RoundRobinCache.Set(rgName, nr, cache.DefaultExpiration)
	nextRoute := routeTable.([]string)[nr.(int)]
	return nextRoute
}

func SetRouteTable(orgID int, rgName string, routeTable []string) {
	tag := fmt.Sprintf("%d-%s", orgID, rgName)
	RoundRobinRouteTable.Set(tag, routeTable, cache.NoExpiration)
}
