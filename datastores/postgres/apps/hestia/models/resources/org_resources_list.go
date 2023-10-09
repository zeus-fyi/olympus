package hestia_compute_resources

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type OrgResourcesGroups struct {
	OrgResourceNodes
}

type OrgResourceDisks struct {
	hestia_autogen_bases.Disks
	hestia_autogen_bases.OrgResources
}

type OrgResourceNodes struct {
	hestia_autogen_bases.Nodes
	hestia_autogen_bases.OrgResources
}

func SelectOrgResourcesNodes(ctx context.Context, orgID int) ([]OrgResourceNodes, error) {
	q := `SELECT ou.org_resource_id, ou.quantity, ou.begin_service, ou.end_service, ou.free_trial, ou.resource_id,
		 n.slug, n.description,n.memory, n.memory_units, n.vcpus, n.disk, n.disk_units,
		  n.price_monthly, n.price_hourly, n.region, n.cloud_provider
		FROM org_resources ou
		JOIN nodes n ON n.resource_id = ou.resource_id
		WHERE org_id = $1 AND ou.end_service IS NULL`

	args := []interface{}{
		orgID,
	}
	// Execute the SQL query
	rows, err := apps.Pg.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the result into a NodesSlice
	var orgResourceNodesSlice []OrgResourceNodes
	for rows.Next() {
		orgResourceNodes := OrgResourceNodes{}

		err = rows.Scan(
			&orgResourceNodes.OrgResourceID,
			&orgResourceNodes.Quantity,
			&orgResourceNodes.BeginService,
			&orgResourceNodes.EndService,
			&orgResourceNodes.FreeTrial,
			&orgResourceNodes.OrgResources.ResourceID,
			&orgResourceNodes.Slug,
			&orgResourceNodes.Description,
			&orgResourceNodes.Memory,
			&orgResourceNodes.MemoryUnits,
			&orgResourceNodes.Vcpus,
			&orgResourceNodes.Disk,
			&orgResourceNodes.DiskUnits,
			&orgResourceNodes.PriceMonthly,
			&orgResourceNodes.PriceHourly,
			&orgResourceNodes.Region,
			&orgResourceNodes.CloudProvider,
		)
		orgResourceNodes.Nodes.ResourceID = orgResourceNodes.OrgResources.ResourceID
		switch orgResourceNodes.CloudProvider {
		case "aws":
			orgResourceNodes.PriceHourly *= 1.30  // Add 30% to the price
			orgResourceNodes.PriceMonthly *= 1.30 // Add 30% to the price
		case "gcp":
			orgResourceNodes.PriceHourly *= 1.30  // Add 30% to the price
			orgResourceNodes.PriceMonthly *= 1.30 // Add 30% to the price
		case "ovh":
			orgResourceNodes.PriceHourly *= 1.20  // Add 20% to the price
			orgResourceNodes.PriceMonthly *= 1.20 // Add 20% to the price
		default:
			orgResourceNodes.PriceHourly *= 1.2  // Add 20% to the price
			orgResourceNodes.PriceMonthly *= 1.2 // Add 20% to the price
		}
		if err != nil {
			return nil, err
		}
		orgResourceNodesSlice = append(orgResourceNodesSlice, orgResourceNodes)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orgResourceNodesSlice, nil
}

func SelectOrgResourcesDisks(ctx context.Context, orgID int) ([]OrgResourceDisks, error) {
	// Build the SQL query
	q := `SELECT ou.quantity, ou.begin_service, ou.end_service, ou.free_trial, ou.resource_id,
		 d.description, d.disk_size, d.disk_units, d.price_monthly, d.price_hourly, d.region, d.cloud_provider
		FROM org_resources ou
		JOIN disks d ON d.resource_id = ou.resource_id
		WHERE org_id = $1`
	args := []interface{}{
		orgID,
	}
	// Execute the SQL query
	rows, err := apps.Pg.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the result into a NodesSlice
	var orgResourceDisksSlice []OrgResourceDisks
	for rows.Next() {
		orgResourceNodes := OrgResourceDisks{}

		err = rows.Scan(
			&orgResourceNodes.Quantity,
			&orgResourceNodes.BeginService,
			&orgResourceNodes.EndService,
			&orgResourceNodes.FreeTrial,
			&orgResourceNodes.Disks.ResourceID,
			&orgResourceNodes.Description,
			&orgResourceNodes.DiskSize,
			&orgResourceNodes.DiskUnits,
			&orgResourceNodes.PriceMonthly,
			&orgResourceNodes.PriceHourly,
			&orgResourceNodes.Region,
			&orgResourceNodes.CloudProvider,
		)
		switch orgResourceNodes.CloudProvider {
		case "aws":
			orgResourceNodes.PriceHourly *= 1.30  // Add 30% to the price
			orgResourceNodes.PriceMonthly *= 1.30 // Add 30% to the price
		case "gcp":
			orgResourceNodes.PriceHourly *= 1.30  // Add 30% to the price
			orgResourceNodes.PriceMonthly *= 1.30 // Add 30% to the price
		case "ovh":
			orgResourceNodes.PriceHourly *= 1.20  // Add 20% to the price
			orgResourceNodes.PriceMonthly *= 1.20 // Add 20% to the price
		default:
			orgResourceNodes.PriceHourly *= 1.2  // Add 20% to the price
			orgResourceNodes.PriceMonthly *= 1.2 // Add 20% to the price
		}
		if err != nil {
			return nil, err
		}
		orgResourceDisksSlice = append(orgResourceDisksSlice, orgResourceNodes)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orgResourceDisksSlice, nil
}

func SelectOrgResourcesDisksAtCloudCtxNs(ctx context.Context, orgID int, cloudCtxNs zeus_common_types.CloudCtxNs) ([]OrgResourceDisks, error) {
	q := `
		WITH cte_get_cloud_ctx AS (
			SELECT cloud_ctx_ns_id 
			FROM topologies_org_cloud_ctx_ns
			WHERE cloud_provider = $2 AND region = $3 AND context = $4 AND namespace = $5 AND org_id = $1
			LIMIT 1
		)
		SELECT ou.org_resource_id, r.resource_id
		FROM org_resources ou
		LEFT JOIN org_resources_cloud_ctx oc ON oc.org_resource_id = ou.org_resource_id
		LEFT JOIN resources r ON r.resource_id = ou.resource_id
		LEFT JOIN disks d ON d.resource_id = r.resource_id
		WHERE org_id = $1 AND r.type = 'disk' AND oc.cloud_ctx_ns_id IN (SELECT cloud_ctx_ns_id FROM cte_get_cloud_ctx)`
	args := []interface{}{
		orgID, cloudCtxNs.CloudProvider, cloudCtxNs.Region, cloudCtxNs.Context, cloudCtxNs.Namespace,
	}
	// Execute the SQL query
	rows, err := apps.Pg.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the result into a NodesSlice
	var orgResourceDisksSlice []OrgResourceDisks
	for rows.Next() {
		orgResourceDisks := OrgResourceDisks{}

		err = rows.Scan(
			&orgResourceDisks.OrgResources.OrgResourceID,
			&orgResourceDisks.OrgResources.ResourceID,
		)
		orgResourceDisks.Disks.ResourceID = orgResourceDisks.OrgResources.ResourceID
		if err != nil {
			return nil, err
		}
		orgResourceDisksSlice = append(orgResourceDisksSlice, orgResourceDisks)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orgResourceDisksSlice, nil
}
