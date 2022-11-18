package zeus

import (
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type K8sRequest struct {
	Kns autok8s_core.CloudCtxNs
}

var K8Util autok8s_core.K8Util
