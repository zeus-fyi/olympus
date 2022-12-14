package containers

import (
	"fmt"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodTemplateSpec) insertContainerSecurityCtx(c containers.Container, contVmsSubCTE, contVmsRelationshipsSubCTE *sql_query_templates.SubCTE) {
	ts := chronos.Chronos{}
	securityCtxID := ts.UnixTimeStampNow()
	contVmsSubCTE.AddValues(securityCtxID, c.SecurityContext.SecurityContextKeyValues)
	contVmsRelationshipsSubCTE.AddValues(securityCtxID, c.GetContainerID())
	return
}

func CreateSecCtxCTEs() (sql_query_templates.SubCTE, sql_query_templates.SubCTE) {
	// env vars
	ts := chronos.Chronos{}
	containerSecurityCtx := autogen_bases.ContainerSecurityContext{}
	containerSecurityCtxSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_security_ctx_%d", ts.UnixTimeStampNow()))
	containerSecurityCtxSubCTE.TableName = containerSecurityCtx.GetTableName()
	containerSecurityCtxSubCTE.Columns = []string{"container_security_context_id", "security_context_key_values"}

	containersSecurityRelationCtx := autogen_bases.ContainersSecurityContext{}
	containersSecurityCtxRelationSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_security_ctx_relationships_%d", ts.UnixTimeStampNow()))
	containersSecurityCtxRelationSubCTE.TableName = containersSecurityRelationCtx.GetTableName()
	containersSecurityCtxRelationSubCTE.Columns = []string{"container_security_context_id", "container_id"}
	return containerSecurityCtxSubCTE, containersSecurityCtxRelationSubCTE
}
