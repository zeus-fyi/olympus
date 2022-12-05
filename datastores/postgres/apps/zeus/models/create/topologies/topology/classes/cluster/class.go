package create_clusters

//type ClusterClass struct {
//	ClassName       string            `json:"className"`
//	OrgUser         org_users.OrgUser `json:"orgUser"`
//	TopologyClassID int               `json:"topologyClassID"`
//	SkeletonBaseID []int             `json:"skeletonBaseIDs"`
//}
//
//func (c *Cluster) AddClusterClassType(cc ClusterClass) ClusterClass {
//	c.ClusterClass = cc
//	return c.ClusterClass
//}
//
//func InsertClusterClassQ() sql_query_templates.QueryParams {
//	q := sql_query_templates.QueryParams{}
//	q.QueryName = "InsertClusterClassDefinition"
//	q.RawQuery = `WITH cte_insert_cluster_class AS (
//					SELECT 1
//				  )
//				  INSERT INTO topology_base_components (org_id, topology_class_type_id, topology_system_component_id, topology_base_name)
//			      VALUES ($1, $2, $3, $4)
//			      RETURNING topology_base_component_id`
//	return q
//}
//
//func InsertClusterClass(ctx context.Context, cc ClusterClass) error {
//	q := InsertClusterClassQ()
//	log.Debug().Interface("InsertClusterClass:", q.LogHeader(Sn))
//
//	return misc.ReturnIfErr(nil, q.LogHeader(Sn))
//}
