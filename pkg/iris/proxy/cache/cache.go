package iris_cache

import (
	"github.com/oleiade/lane/v2"
	"github.com/puzpuzpuz/xsync/v2"
	iris_proxy "github.com/zeus-fyi/olympus/pkg/iris/proxy"
)

var RouteCache = xsync.NewMapOf[string]()

func NewRequestQueue() *lane.Queue[iris_proxy.ProxyRequest] {
	return lane.NewQueue[iris_proxy.ProxyRequest]()
}

type ProxyRequestPriorityQueue struct {
	Type          string
	PriorityQueue *lane.Queue[iris_proxy.ProxyRequest]
}
