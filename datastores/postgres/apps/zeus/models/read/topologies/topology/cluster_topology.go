package read_topology

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
)

func selectClusterTopologiesQ(cte *sql_query_templates.CTE, orgID int, clusterName string, clusterSkeletonBases []string) sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	q.QueryName = "SelectInfraTopologyQuery"
	cte.Params = append(cte.Params, orgID)
	cte.Params = append(cte.Params, clusterName)
	sbs := make([]string, len(clusterSkeletonBases))

	for i, sb := range clusterSkeletonBases {
		sbs[i] = sb
	}
	cte.Params = append(cte.Params, pq.Array(sbs))

	query := `  SELECT topology_base_name, topology_skeleton_base_name, MAX(topology_id)
				FROM topology_system_components ts
				INNER JOIN topology_base_components tb ON tb.topology_system_component_id = ts.topology_system_component_id
				INNER JOIN topology_skeleton_base_components sb ON tb.topology_base_component_id = sb.topology_base_component_id
				INNER JOIN topology_infrastructure_components ti ON ti.topology_skeleton_base_id = sb.topology_skeleton_base_id
				WHERE (ts.org_id = $1 AND ts.topology_system_component_name = $2) AND sb.topology_skeleton_base_name = ANY($3::text[])
				GROUP BY topology_base_name, topology_skeleton_base_name`

	q.RawQuery = query
	return q
}

func selectAllClusterTopologiesQByID(cte *sql_query_templates.CTE, orgID int, appID int) sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	q.QueryName = "SelectInfraTopologyQuery"
	cond := "WHERE (ts.org_id = $1 AND ts.topology_system_component_id = $2)"
	cte.Params = append(cte.Params, orgID)
	cte.Params = append(cte.Params, appID)

	query := fmt.Sprintf(`SELECT ts.topology_system_component_name, tb.topology_base_name, sb.topology_skeleton_base_name, MAX(topology_id)
				FROM topology_system_components ts
				INNER JOIN topology_base_components tb ON tb.topology_system_component_id = ts.topology_system_component_id
				INNER JOIN topology_skeleton_base_components sb ON tb.topology_base_component_id = sb.topology_base_component_id
				INNER JOIN topology_infrastructure_components ti ON ti.topology_skeleton_base_id = sb.topology_skeleton_base_id
				%s
				GROUP BY ts.topology_system_component_name, tb.topology_base_name, sb.topology_skeleton_base_name`, cond)

	q.RawQuery = query
	return q
}

func selectAllClusterTopologiesQByName(cte *sql_query_templates.CTE, orgID int, appName string) sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	q.QueryName = "SelectInfraTopologyQuery"
	cond := "WHERE (ts.org_id = $1 AND ts.topology_system_component_name = $2)"
	cte.Params = append(cte.Params, orgID)
	cte.Params = append(cte.Params, appName)

	query := fmt.Sprintf(`SELECT ts.topology_system_component_name, tb.topology_base_name, sb.topology_skeleton_base_name, MAX(topology_id)
				FROM topology_system_components ts
				INNER JOIN topology_base_components tb ON tb.topology_system_component_id = ts.topology_system_component_id
				INNER JOIN topology_skeleton_base_components sb ON tb.topology_base_component_id = sb.topology_base_component_id
				INNER JOIN topology_infrastructure_components ti ON ti.topology_skeleton_base_id = sb.topology_skeleton_base_id
				%s
				GROUP BY ts.topology_system_component_name, tb.topology_base_name, sb.topology_skeleton_base_name`, cond)

	q.RawQuery = query
	return q
}

type ClusterTopology struct {
	ClusterClassName string              `json:"clusterName"`
	Topologies       []ClusterTopologies `json:"topologies"`
}

func (c ClusterTopology) GetTopologyIDs() []int {
	tmp := make([]int, len(c.Topologies))
	for i, ct := range c.Topologies {
		tmp[i] = ct.TopologyID
	}
	return tmp
}

func (c ClusterTopology) CheckForChoreographyOption() bool {
	for _, ct := range c.Topologies {
		if ct.SkeletonBaseName == "choreography" || ct.SkeletonBaseName == "hydraChoreography" {
			return true
		}
	}
	return false
}

type ClusterTopologies struct {
	TopologyID        int    `json:"topologyID"`
	ComponentBaseName string `json:"componentBaseName"`
	SkeletonBaseName  string `json:"skeletonBaseName"`
	Tag               string `json:"tag"`
}

func SelectClusterTopologyFiltered(ctx context.Context, orgID int, clusterName string, clusterSkeletonBases []string, filter map[string]map[string]bool) (ClusterTopology, error) {
	res, err := SelectClusterTopology(ctx, orgID, clusterName, clusterSkeletonBases)
	if err != nil {
		return res, err
	}
	filteredTop := []ClusterTopologies{}
	for _, val := range res.Topologies {
		_, ok := filter[val.ComponentBaseName][val.SkeletonBaseName]
		if ok {
			filteredTop = append(filteredTop, val)
		}
	}
	res.Topologies = filteredTop
	return res, nil
}

func SelectClusterTopology(ctx context.Context, orgID int, clusterName string, clusterSkeletonBases []string) (ClusterTopology, error) {
	cte := sql_query_templates.CTE{}
	cl := ClusterTopology{ClusterClassName: clusterName}
	cl.Topologies = []ClusterTopologies{}
	q := selectClusterTopologiesQ(&cte, orgID, clusterName, clusterSkeletonBases)
	log.Debug().Interface("SelectTopologiesMetadata:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, cte.Params...)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return cl, err
	}
	defer rows.Close()
	for rows.Next() {
		ct := ClusterTopologies{Tag: "latest"}
		rowErr := rows.Scan(
			&ct.ComponentBaseName, &ct.SkeletonBaseName, &ct.TopologyID,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return cl, rowErr
		}
		cl.Topologies = append(cl.Topologies, ct)
	}
	return cl, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func SelectAppTopologyByID(ctx context.Context, orgID int, appID int) (zeus_cluster_config_drivers.ClusterDefinition, error) {
	cte := sql_query_templates.CTE{}
	cl := zeus_cluster_config_drivers.ClusterDefinition{
		ComponentBases: make(map[string]zeus_cluster_config_drivers.ComponentBaseDefinition),
	}
	q := selectAllClusterTopologiesQByID(&cte, orgID, appID)
	log.Debug().Interface("SelectAppTopologyByID:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, cte.Params...)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return cl, err
	}
	defer rows.Close()
	for rows.Next() {
		var componentBaseName, skeletonBaseName string
		var topologyID int
		rowErr := rows.Scan(
			&cl.ClusterClassName, &componentBaseName, &skeletonBaseName, &topologyID,
		)
		if _, ok := cl.ComponentBases[componentBaseName]; !ok {
			cl.ComponentBases[componentBaseName] = zeus_cluster_config_drivers.ComponentBaseDefinition{
				SkeletonBases: make(map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition),
			}
		}
		if _, ok := cl.ComponentBases[componentBaseName].SkeletonBases[skeletonBaseName]; !ok {
			cl.ComponentBases[componentBaseName].SkeletonBases[skeletonBaseName] = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{}
		}
		tmp := cl.ComponentBases[componentBaseName].SkeletonBases[skeletonBaseName]
		tmp.TopologyID = topologyID
		cl.ComponentBases[componentBaseName].SkeletonBases[skeletonBaseName] = tmp
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return cl, rowErr
		}
	}
	return cl, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func SelectAppTopologyByName(ctx context.Context, orgID int, appName string) (zeus_cluster_config_drivers.ClusterDefinition, error) {
	cte := sql_query_templates.CTE{}
	cl := zeus_cluster_config_drivers.ClusterDefinition{
		ComponentBases: make(map[string]zeus_cluster_config_drivers.ComponentBaseDefinition),
	}
	q := selectAllClusterTopologiesQByName(&cte, orgID, appName)
	log.Debug().Interface("SelectAppTopologyByID:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, cte.Params...)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return cl, err
	}
	defer rows.Close()
	for rows.Next() {
		var componentBaseName, skeletonBaseName string
		var topologyID int
		rowErr := rows.Scan(
			&cl.ClusterClassName, &componentBaseName, &skeletonBaseName, &topologyID,
		)
		if _, ok := cl.ComponentBases[componentBaseName]; !ok {
			cl.ComponentBases[componentBaseName] = zeus_cluster_config_drivers.ComponentBaseDefinition{
				SkeletonBases: make(map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition),
			}
		}
		if _, ok := cl.ComponentBases[componentBaseName].SkeletonBases[skeletonBaseName]; !ok {
			cl.ComponentBases[componentBaseName].SkeletonBases[skeletonBaseName] = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{}
		}
		tmp := cl.ComponentBases[componentBaseName].SkeletonBases[skeletonBaseName]
		tmp.TopologyID = topologyID
		cl.ComponentBases[componentBaseName].SkeletonBases[skeletonBaseName] = tmp
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return cl, rowErr
		}
	}
	return cl, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
