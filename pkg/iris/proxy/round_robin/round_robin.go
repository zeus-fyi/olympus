package iris_round_robin

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

var (
	RoundRobinCache      = cache.New(cache.NoExpiration, cache.NoExpiration)
	RoundRobinRouteTable = cache.New(cache.NoExpiration, cache.NoExpiration)
	NotFound             = fmt.Errorf("routing group not found")
)

func orgRouteTag(orgID int, rgName string) string {
	return fmt.Sprintf("%d-%s", orgID, rgName)
}

func GetNextRoute(orgID int, rgName string) (string, error) {
	// TODO IrisCache here
	tag := orgRouteTag(orgID, rgName)
	routeTable, ok := RoundRobinRouteTable.Get(tag)
	if !ok {
		return "", fmt.Errorf("no route table found")
	}
	nr, ok := RoundRobinCache.Get(rgName)
	if !ok {
		return "", fmt.Errorf("no route group found")
	}
	nr = (nr.(int) + 1) % len(routeTable.([]string))
	RoundRobinCache.Set(rgName, nr, cache.DefaultExpiration)
	nextRoute := routeTable.([]string)[nr.(int)]
	return nextRoute, nil
}

//func SetRouteTable(orgID int, rgName string, routeTable []string) {
//	// TODO IrisCache here
//	tag := orgRouteTag(orgID, rgName)
//	RoundRobinCache.Set(rgName, 0, cache.DefaultExpiration)
//	RoundRobinRouteTable.Set(tag, routeTable, cache.NoExpiration)
//}
//
//func InitRoutingTables(ctx context.Context) {
//	// TODO IrisCache here
//	ot, err := iris_models.SelectAllOrgRoutes(ctx)
//	if err != nil {
//		panic(err)
//	}
//	for orgID, og := range ot.Map {
//		for rgName, routeTable := range og {
//			SetRouteTable(orgID, rgName, routeTable)
//		}
//	}
//}
