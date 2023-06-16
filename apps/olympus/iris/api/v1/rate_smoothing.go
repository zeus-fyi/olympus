package v1_iris

import (
	"context"
	"fmt"

	"github.com/oleiade/lane/v2"
	"github.com/rs/zerolog/log"
)

func NewRequestQueue() *lane.Queue[ProxyRequest] {
	return lane.NewQueue[ProxyRequest]()
}

type ProxyRequestPriorityQueue struct {
	Type          string
	PriorityQueue *lane.Queue[ProxyRequest]
}

func (q *ProxyRequestPriorityQueue) ProcessQueue(ctx context.Context) {
	ql := q.PriorityQueue.Size()
	if ql == 0 {
		return
	}
	log.Info().Str("signingType", q.Type).Msg(fmt.Sprintf("queue length: %d", ql))
	for {
		ql = q.PriorityQueue.Size()
		if ql == 0 {
			break
		}
		_, qOk := q.PriorityQueue.Dequeue()
		if !qOk {
			continue
		}
	}
}
