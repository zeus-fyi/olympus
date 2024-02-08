package topology_auths

import (
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

var (
	K8Util  autok8s_core.K8Util
	KeysCfg auth_startup.AuthConfig
)
