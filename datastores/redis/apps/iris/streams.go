package iris_redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	maxStreamLen         = 100
	EthMempoolStreamName = "eth-mempool-mainnet"
)

// CreateOrAddToStream initializes a new stream topic if it doesn't already exist.
func (m *IrisCache) CreateOrAddToStream(ctx context.Context, streamName string, payload map[string]interface{}) error {
	// Use the XADD command with the special ID '*' to add an item and create the stream if it doesn't exist.
	_, err := m.Writer.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: payload,
		MaxLen: maxStreamLen,
	}).Result()
	if err != nil {
		log.Err(err).Msgf("error creating redis stream: %s\n", err.Error())
		return err
	}
	return nil
}

// Stream fetches data from the given stream.
func (m *IrisCache) Stream(ctx context.Context, streamName string, lastID string) ([]redis.XStream, error) {
	// Use the XREAD command to read from the stream.
	messages, err := m.Reader.XRead(ctx, &redis.XReadArgs{
		Streams: []string{streamName, lastID},
		Count:   maxStreamLen / 2, // Or any number you prefer to limit the fetched messages.
		Block:   0,
	}).Result()
	if err != nil {
		log.Err(err).Msgf("error reading redis stream: %s\n", err.Error())
		return nil, err
	}
	if len(messages) == 0 {
		return nil, nil
	}
	return messages, nil
}
