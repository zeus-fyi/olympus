package iris_redis

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func orgRouteTag(orgID int, rgName string) string {
	return fmt.Sprintf("%d-%s", orgID, rgName)
}

func (m *IrisCache) GetNextRoute(ctx context.Context, orgID int, rgName string) (string, error) {
	// Generate the key
	tag := orgRouteTag(orgID, rgName)

	// Use Redis transaction (pipeline) to perform check existence and get operation atomically
	pipe := m.Reader.TxPipeline()

	// Check if the key exists
	existsCmd := pipe.Exists(ctx, tag)

	// Pop the endpoint from the head of the list
	endpointCmd := pipe.LPop(ctx, tag)

	// Push the popped endpoint back to the tail, ensuring round-robin rotation
	pipe.RPush(ctx, tag, endpointCmd.Val())

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("error during pipeline execution for key: %s\n", tag)
		return "", err
	}

	// Check whether the key exists
	if existsCmd.Val() <= 0 {
		return "", fmt.Errorf("key doesn't exist: %s", tag)
	}

	// Return the endpoint
	return endpointCmd.Val(), nil
}

func (m *IrisCache) AddOrUpdateOrgRoutingGroup(ctx context.Context, orgID int, rgName string, routes []string) error {
	// Generate the key
	tag := orgRouteTag(orgID, rgName)

	// Start a new transaction
	pipe := m.Writer.TxPipeline()

	// Remove the old key if it exists
	pipe.Del(ctx, tag)

	// Add each route to the list
	for _, route := range routes {
		pipe.RPush(ctx, tag, route)
	}

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		fmt.Printf("error updating routing group: %s\n", tag)
		return err
	}

	return nil
}

func (m *IrisCache) initRoutingTables(ctx context.Context) error {
	ot, err := iris_models.SelectAllOrgRoutes(ctx)
	if err != nil {
		return err
	}
	for orgID, og := range ot.Map {
		for rgName, routes := range og {
			err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
			if err != nil {
				log.Err(err).Msg("InitRoutingTables")
				return err
			}
		}
	}
	return nil
}

func (m *IrisCache) initRoutingTablesForOrg(ctx context.Context, orgID int) error {
	ot, err := iris_models.SelectAllOrgRoutesByOrg(ctx, orgID)
	if err != nil {
		return err
	}
	for _, og := range ot.Map {
		for rgName, routes := range og {
			err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
			if err != nil {
				log.Err(err).Msg("InitRoutingTablesForOrg")
				return err
			}
		}
	}
	return nil
}

func (m *IrisCache) refreshRoutingTablesForOrgGroup(ctx context.Context, orgID int, groupName string) error {
	ot, err := iris_models.SelectOrgRoutesByOrgAndGroupName(ctx, orgID, groupName)
	if err != nil {
		return err
	}
	for _, og := range ot.Map {
		for rgName, routes := range og {
			err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
			if err != nil {
				log.Err(err).Msg("InitRoutingTablesForOrgGroup")
				return err
			}
		}
	}
	return nil
}
