package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodTemplateSpec) insertContainerSecurityCtx(m map[string]containers.Container, imageID string, contVmsSubCTE, contVmsRelationshipsSubCTE *sql_query_templates.SubCTE) {
	c, ok := m[imageID]
	if !ok {
		return
	}
	ts := chronos.Chronos{}
	securityCtxID := ts.UnixTimeStampNow()
	contVmsSubCTE.AddValues(securityCtxID, c.SecurityContext.SecurityContextKeyValues)
	contVmsRelationshipsSubCTE.AddValues(securityCtxID, c.GetContainerID())
	return
}
