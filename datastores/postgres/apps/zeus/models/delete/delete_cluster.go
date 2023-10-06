package delete_cluster

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

var q = `WITH cte_delete_cluster_id AS (
			SELECT topology_system_component_id
			FROM topology_system_components WHERE topology_system_component_name = $2 AND org_id = $1
		), cte_delete_bases_id AS (
			SELECT topology_base_component_id
			FROM topology_base_components tb
			WHERE tb.topology_system_component_id = (SELECT topology_system_component_id FROM cte_delete_cluster_id)
		), cte_delete_sb_id AS (
			SELECT topology_skeleton_base_id
			FROM topology_skeleton_base_components sb 
			WHERE sb.topology_base_component_id IN (SELECT topology_base_component_id FROM cte_delete_bases_id)
		), cte_delete_sb_infra AS (
			DELETE
			FROM topology_infrastructure_components
			WHERE topology_skeleton_base_id IN (SELECT topology_skeleton_base_id FROM cte_delete_sb_id)
		), cte_delete_sb AS (
			DELETE
			FROM topology_skeleton_base_components
			WHERE topology_skeleton_base_id IN (SELECT topology_skeleton_base_id FROM cte_delete_sb_id)
		), cte_delete_cb AS (
			DELETE
			FROM topology_base_components
			WHERE topology_base_component_id IN (SELECT topology_base_component_id FROM cte_delete_bases_id)
		) 
		DELETE
		FROM topology_system_components
		WHERE topology_system_component_id = (SELECT topology_system_component_id FROM cte_delete_cluster_id)`

func DeleteCluster(ctx context.Context, orgID int, name string) error {
	_, err := apps.Pg.Exec(ctx, q, orgID, name)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("DeleteCluster: failed to delete cluster")
		return err
	}
	return nil
}
