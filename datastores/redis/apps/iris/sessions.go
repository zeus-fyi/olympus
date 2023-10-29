package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

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

func (m *IrisCache) AddSessionWithTTL(ctx context.Context, sessionID, route string, ttl time.Duration) error {
	// Start a pipeline using the writer client
	pipe := m.Writer.Pipeline()

	// Add the sessionID with the provided TTL
	setSessionCmd := pipe.Set(ctx, sessionID, route, ttl)

	// Add the route with the provided TTL
	setRouteCmd := pipe.Set(ctx, route, sessionID, ttl)

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("AddSessionWithTTL: %s, %s", sessionID, route)
		return err
	}

	// Check error for setSessionCmd
	err = setSessionCmd.Err()
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("Error in setting sessionID: %s", sessionID)
		return err
	}

	// Check error for setRouteCmd
	err = setRouteCmd.Err()
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("Error in setting route: %s", route)
		return err
	}

	// No errors, return nil
	return nil
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
