package iris_traffic_controllers

import (
	"context"

	"github.com/oleiade/lane/v2"
	iris_proxy "github.com/zeus-fyi/olympus/pkg/iris/proxy"
)

func NewRequestQueue() *lane.Queue[iris_proxy.ProxyRequest] {
	return lane.NewQueue[iris_proxy.ProxyRequest]()
}

type ProxyRequestPriorityQueue struct {
	Type          string
	PriorityQueue *lane.Queue[iris_proxy.ProxyRequest]
}

func (q *ProxyRequestPriorityQueue) ProcessQueue(ctx context.Context) {
	ql := q.PriorityQueue.Size()
	if ql == 0 {
		return
	}
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
