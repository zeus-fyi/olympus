package create_clusters

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "Cluster"

/*
--------------------------------------------------------------------------------------------------------------------------------
linking to make cluster topology

steps

0. deploy table 2 migrate
1. create topology_system_components

topology_system_components (check class_id >= 3)
 ethereumBeacon, 4
 execClient, 3
 consensusClient, 3

// skeleton bases
 geth, 2
	-> connect to infra components
	-> TODO connect to config components
 lighthouse, 2
	-> connect to infra components
	-> TODO connect to config components

deploy latest cluster
	-> uses system component id
		-> queries all base components
			-> queries latest (others todo) skeleton base component ids
				-> queries all infra (+config todo) to get topology_ids & chart packages related to deploy
					-> deploys all chart packages to namespace

2. deploy beacon using a cluster id
	-> where (skeletonBases = geth, lighthouse), version = latest

With cte_insert_cluster_topology AS (
	INSERT INTO topology_class_types (topology_class_type_id
)
*/

// InsertCluster should probably be using pre-stored bases and just really managing the kns creation
// and topology dependencies being deployed, eg taking any config -> apply to infra -> deploy
func (c *Cluster) InsertCluster(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))

	return misc.ReturnIfErr(nil, q.LogHeader(Sn))
}

/* TODO

--------------------------------------------------------------------------------------------------------------------------------
changing config, linking infra to config

topology_infrastructure_components
    "topology_infrastructure_component_id" int8 DEFAULT next_id(),
    "topology_id" int8 NOT NULL REFERENCES topologies(topology_id),
    "chart_package_id" int8 NOT NULL REFERENCES chart_packages(chart_package_id)

not used yet...

topology_configuration_class
    "topology_configuration_class_id" int8 DEFAULT next_id(),
    "topology_infrastructure_component_id" int8 NOT NULL REFERENCES topology_infrastructure_components(topology_infrastructure_component_id)

topology_configuration_child_values_overrides
    "topology_configuration_child_values_override_id" int8 DEFAULT next_id(),
    "topology_configuration_class_id" int8 NOT NULL REFERENCES topology_configuration_class(topology_configuration_class_id),
    "chart_subcomponent_child_values_id" int8 NOT NULL REFERENCES chart_subcomponents_child_values(chart_subcomponent_child_values_id),
    "chart_subcomponent_override_value" text NOT NULL
--------------------------------------------------------------------------------------------------------------------------------

*/
