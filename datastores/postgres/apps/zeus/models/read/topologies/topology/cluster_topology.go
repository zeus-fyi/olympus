package read_topology

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func selectClusterTopologiesQ(cte *sql_query_templates.CTE, orgID int, clusterName string, clusterSkeletonBases []string) sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	q.QueryName = "SelectInfraTopologyQuery"
	cond := "WHERE (ts.org_id = $1 AND ts.topology_system_component_name = $2)"
	cte.Params = append(cte.Params, orgID)
	cte.Params = append(cte.Params, clusterName)
	for i, skb := range clusterSkeletonBases {
		if i == 0 {
			cond += " AND ("
		}
		cte.Params = append(cte.Params, skb)
		cond += fmt.Sprintf("topology_skeleton_base_name = $%d", len(cte.Params))
		if i != len(clusterSkeletonBases)-1 {
			cond += " OR "
		}
		if i == len(clusterSkeletonBases)-1 {
			cond += ")"
		}
	}

	query := fmt.Sprintf(`SELECT topology_skeleton_base_name, MAX(topology_id)
				FROM topology_system_components ts
				INNER JOIN topology_base_components tb ON tb.topology_system_component_id = ts.topology_system_component_id
				INNER JOIN topology_skeleton_base_components sb ON tb.topology_base_component_id = sb.topology_base_component_id
				INNER JOIN topology_infrastructure_components ti ON ti.topology_skeleton_base_id = sb.topology_skeleton_base_id
				%s
				GROUP BY topology_skeleton_base_name`, cond)

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
	TopologyID       int    `json:"topologyID"`
	SkeletonBaseName string `json:"skeletonBaseName"`
	Tag              string `json:"tag"`
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
			&ct.SkeletonBaseName, &ct.TopologyID,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return cl, rowErr
		}
		cl.Topologies = append(cl.Topologies, ct)
	}
	return cl, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
