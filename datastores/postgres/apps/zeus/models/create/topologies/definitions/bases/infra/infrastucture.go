package create_infra

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/infra"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type InfraBaseTopology struct {
	infra.InfraBaseTopology
	packages.Packages

	Tag               string `json:"tag,omitempty"`
	ClusterClassName  string `json:"clusterClassName,omitempty"`
	ComponentBaseName string `json:"componentBaseName,omitempty"`
	SkeletonBaseName  string `json:"skeletonBaseName,omitempty"`
	TopologyClassID   int    `json:"topologyClassID,omitempty"`
}

func NewCreateInfrastructure() InfraBaseTopology {
	pkg := packages.NewPackageInsert()
	infc := infra.NewInfrastructureBaseTopology()
	ibc := InfraBaseTopology{infc, pkg, "latest-internal", "", "", "", 0}
	return ibc
}

func InsertInfraBeaconCopy(ctx context.Context, ou org_users.OrgUser) error {
	var q sql_query_templates.QueryParams
	_, err := apps.Pg.Exec(ctx, COPY_BEACON, ou.OrgID, ou.UserID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

var COPY_BEACON = `WITH cte_insert_max AS (
			SELECT MAX(tc.chart_package_id) AS chart_package_id, cp.chart_name AS chart_name, (SELECT next_id() - topology_skeleton_base_id) AS tid, topology_skeleton_base_id AS topology_skeleton_base_id, 'copy' AS tag
			FROM topology_infrastructure_components tc
			JOIN chart_packages cp ON cp.chart_package_id = tc.chart_package_id
			WHERE topology_skeleton_base_id = 1670472660668679974 OR topology_skeleton_base_id = 1670472660615421388
			GROUP BY topology_skeleton_base_id, chart_name
		), cte_it AS (
			INSERT INTO topologies (name, topology_id) 
			SELECT chart_name, tid
			FROM cte_insert_max
		), cte_orguser_top AS (
			INSERT INTO org_users_topologies (org_id, user_id, topology_id) 
			SELECT $1, $2, tid
			FROM cte_insert_max
		)  INSERT INTO topology_infrastructure_components (topology_infrastructure_component_id, topology_id, chart_package_id, tag, topology_skeleton_base_id) 
		   SELECT tid, tid, chart_package_id, tag, topology_skeleton_base_id
 		   FROM cte_insert_max
`
