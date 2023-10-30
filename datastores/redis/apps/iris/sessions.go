package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

//type ServerlessRoute struct {
//	OrgID      int
//	SessionID  string
//	TableName  string
//	TableRoute string
//}

func (m *IrisCache) AddSessionWithTTL(ctx context.Context, orgID int, sessionID, tableRoute, serverlessRoutesTable string, ttl time.Duration) error {
	// Start a pipeline using the writer client
	orgSessionID := getOrgSessionIDKey(orgID, sessionID)
	pipe := m.Writer.Pipeline()
	/*
		user -> sessions map
			getOrgSessionIDKey(orgID, sessionID) -> route

		    1. add session -> tableRoute
		    2. add session -> serverlessRoute
	*/
	//serviceTable := getGlobalServerlessTableKey(serverlessRoutesTable)

	//minElemCmd := pipe.ZRangeWithScores(ctx, serviceTable, 0, 0)
	//setSessionCmd := pipe.Set(ctx, orgSessionID, serverlessRoute, ttl)
	setRouteCmd := pipe.Set(ctx, orgSessionID, tableRoute, ttl)

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		// If there's an error, log it and return it
		log.Err(err).Msgf("AddSessionWithTTL: %s, %s", orgSessionID, tableRoute)
		return err
	}

	// Check error for setSessionCmd
	//err = setSessionCmd.Err()
	//if err != nil {
	//	// If there's an error, log it and return it
	//	log.Err(err).Msgf("AddSessionWithTTL: %s, %s, %s", orgSessionID,  tableRoute)
	//	return err
	//}

	// Check error for setRouteCmd
	err = setRouteCmd.Err()
	if err != nil {
		// If there's an error, log it and return it
		log.Err(err).Msgf("AddSessionWithTTL: %s, %s", orgSessionID, tableRoute)
		return err
	}
	// No errors, return nil
	return nil
}

func (m *IrisCache) DoesSessionIDExist(ctx context.Context, key string) (bool, error) {
	// Check if the key exists
	existsCmd := m.Reader.Exists(ctx, key)

	// Check error
	err := existsCmd.Err()
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("KeyExists: %s", key)
		return false, err
	}

	// Return whether the key exists
	return existsCmd.Val() > 0, nil
}

func (m *IrisCache) GetAndUpdateLatestSessionCacheTTLIfExists(ctx context.Context, sessionID string, ttl time.Duration) (string, error) {
	// Start a pipeline using the writer client
	pipe := m.Writer.Pipeline()

	// Perform the Exists command
	existsCmd := pipe.Exists(ctx, sessionID)

	// Get the value
	getCmd := pipe.Get(ctx, sessionID)

	// Update the TTL
	pipe.Expire(ctx, sessionID, ttl)

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("GetAndUpdateLatestSessionCacheTTLIfExists: %s", sessionID)
		return "", err
	}

	// Check if the key exists
	if existsCmd.Val() == 0 {
		return "", redis.Nil
	}

	// Return the value
	return getCmd.Val(), nil
}

func (m *IrisCache) DeleteSessionCacheIfExists(ctx context.Context, sessionID string) (int64, error) {
	// Delete the key
	delCmd := m.Writer.Del(ctx, sessionID)

	// Check error
	err := delCmd.Err()
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("DeleteSessionCacheIfExists: %s", sessionID)
		return 0, err
	}

	// Return the number of keys deleted
	return delCmd.Val(), nil
}

func (m *IrisCache) DeleteSessionRouteCacheIfExists(ctx context.Context, route string) (int64, error) {
	// Delete the key
	delCmd := m.Writer.Del(ctx, route)

	// Check error
	err := delCmd.Err()
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("DeleteSessionRouteCacheIfExists: %s", route)
		return 0, err
	}

	// Return the number of keys deleted
	return delCmd.Val(), nil
}
