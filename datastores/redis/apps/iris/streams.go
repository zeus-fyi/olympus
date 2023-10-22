package iris_redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

// CreateOrAddToStream initializes a new stream topic if it doesn't already exist.
func (m *IrisCache) CreateOrAddToStream(ctx context.Context, streamName string, payload map[string]interface{}) error {
	// Use the XADD command with the special ID '*' to add an item and create the stream if it doesn't exist.
	_, err := m.Writer.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: payload,
		MaxLen: 2,
	}).Result()
	if err != nil {
		fmt.Printf("error during stream creation: %s\n", err.Error())
		return err
	}
	return nil
}

// Stream fetches data from the given stream.
func (m *IrisCache) Stream(ctx context.Context, streamName string, lastID string) ([]redis.XStream, error) {
	// Use the XREAD command to read from the stream.
	messages, err := m.Reader.XRead(ctx, &redis.XReadArgs{
		Streams: []string{streamName, lastID},
		Count:   10, // Or any number you prefer to limit the fetched messages.
	}).Result()
	if err != nil {
		fmt.Printf("error during stream read: %s\n", err.Error())
		return nil, err
	}
	if len(messages) == 0 {
		return nil, nil
	}
	return messages, nil
}
