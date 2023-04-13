package hestia_compute_resources

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
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
	// Build the SQL query
	q := `SELECT ou.quantity, ou.begin_service, ou.end_service, ou.free_trial, ou.resource_id,
		 n.slug, n.description,n.memory, n.memory_units, n.vcpus, n.disk, n.disk_units,
		  n.price_monthly, n.price_hourly, n.region, n.cloud_provider
		FROM org_resources ou
		JOIN nodes n ON n.resource_id = ou.resource_id
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
	var orgResourceNodesSlice []OrgResourceNodes
	for rows.Next() {
		orgResourceNodes := OrgResourceNodes{}

		err = rows.Scan(
			&orgResourceNodes.Quantity,
			&orgResourceNodes.BeginService,
			&orgResourceNodes.EndService,
			&orgResourceNodes.FreeTrial,
			&orgResourceNodes.Nodes.ResourceID,
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
		orgResourceNodes.PriceHourly *= 1.1  // Add 10% to the price
		orgResourceNodes.PriceMonthly *= 1.1 // Add 10% to the price

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
		orgResourceNodes.PriceHourly *= 1.1  // Add 10% to the price
		orgResourceNodes.PriceMonthly *= 1.1 // Add 10% to the price

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
