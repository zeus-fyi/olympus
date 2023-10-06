package read_topology

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func SelectOrgAppsQ() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "SelectOrgClusters"
	q.RawQuery = `SELECT topology_system_component_id, topology_class_type_id, topology_system_component_name
				  FROM topology_system_components
				  WHERE org_id = $1
				  ORDER BY topology_system_component_name ASC
					`
	return q
}

func SelectOrgApps(ctx context.Context, orgID int) (autogen_bases.TopologySystemComponentsSlice, error) {
	q := SelectOrgAppsQ()
	log.Debug().Interface("SelectOrgClusters:", q.LogHeader(Sn))

	appSlice := autogen_bases.TopologySystemComponentsSlice{}
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return appSlice, err
	}
	defer rows.Close()
	for rows.Next() {
		cluster := autogen_bases.TopologySystemComponents{}
		rowErr := rows.Scan(
			&cluster.TopologySystemComponentID, &cluster.TopologyClassTypeID, &cluster.TopologySystemComponentName,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return appSlice, rowErr
		}
		appSlice = append(appSlice, cluster)
	}

	return appSlice, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

// This is just a stub, until matrix is implemented
func SelectPublicMatrixAppsQ() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "SelectMatrixApps"
	q.RawQuery = `SELECT topology_system_component_id, topology_class_type_id, topology_system_component_name
				  FROM topology_system_components
				  WHERE org_id = $1 AND topology_system_component_name LIKE 'sui-%' 
				  ORDER BY topology_system_component_name ASC
					`
	return q
}

// This is just a stub, until matrix is implemented
func SelectPublicMatrixApps(ctx context.Context, orgID int) (autogen_bases.TopologySystemComponentsSlice, error) {
	q := SelectPublicMatrixAppsQ()
	log.Debug().Interface("SelectMatrixApps:", q.LogHeader(Sn))

	appSlice := autogen_bases.TopologySystemComponentsSlice{}
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return appSlice, err
	}
	defer rows.Close()
	for rows.Next() {
		cluster := autogen_bases.TopologySystemComponents{}
		rowErr := rows.Scan(
			&cluster.TopologySystemComponentID, &cluster.TopologyClassTypeID, &cluster.TopologySystemComponentName,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return appSlice, rowErr
		}
		appSlice = append(appSlice, cluster)
	}

	return appSlice, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
