package iris_redis

import (
	"context"
	"fmt"
)

func (r *IrisRedisTestSuite) TestRedisStreams() {
	streamData := map[string]interface{}{
		"foo":  "bar",
		"baz":  "qux",
		"quux": "quuz",
	}

	for k, v := range streamData {
		m := map[string]interface{}{
			k: v,
		}
		err := IrisRedisClient.CreateOrAddToStream(context.Background(), "test-stream", m)
		r.NoError(err)
	}

	messages, err := IrisRedisClient.Stream(context.Background(), "test-stream", "0")
	r.NoError(err)

	for _, msg := range messages {
		fmt.Println(msg.Stream)
		for _, event := range msg.Messages {
			fmt.Println(event.ID)
			fmt.Println(event.Values)
		}
	}
}
