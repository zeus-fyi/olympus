package iris_redis

import (
	"context"
	"fmt"

	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (m *IrisCache) GetBroadcastRoutes(ctx context.Context, orgID int, rgName string) ([]iris_models.RouteInfo, error) {
	// Use Redis pipeline to perform operations
	pipe := m.Writer.Pipeline()

	// Generate the route key
	routeKey := orgRouteTag(orgID, rgName)

	// Get all elements from the list
	endpointsCmd := pipe.LRange(ctx, routeKey, 0, -1)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("error during pipeline execution: %s\n", err.Error())
		return nil, err
	}

	// Convert the endpoints to the desired type
	var routes []iris_models.RouteInfo
	for _, endpoint := range endpointsCmd.Val() {
		routeInfo := iris_models.RouteInfo{
			RoutePath: endpoint,
			Referers:  nil,
		}
		routes = append(routes, routeInfo)
	}

	return routes, nil
}
