package hestia_compute_resources

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"k8s.io/apimachinery/pkg/api/resource"
)

type NodeFilter struct {
	CloudProvider string                 `json:"cloudProvider"`
	Region        string                 `json:"region"`
	ResourceSums  zeus_core.ResourceSums `json:"resourceSums"`
}

func SelectNodes(ctx context.Context, nf NodeFilter) (hestia_autogen_bases.NodesSlice, error) {
	if nf.ResourceSums.MemRequests == "" {
		nf.ResourceSums.MemRequests = "0"
	}
	memRequests, err := resource.ParseQuantity(nf.ResourceSums.MemRequests)
	if err != nil {
		return nil, err
	}
	if nf.ResourceSums.CpuRequests == "" {
		nf.ResourceSums.CpuRequests = "0"
	}
	cpuRequests, err := resource.ParseQuantity(nf.ResourceSums.CpuRequests)
	if err != nil {
		return nil, err
	}
	// Convert to MegaBytes and vCores
	memRequestsMegaBytes := memRequests.Value() / (1024 * 1024)
	cpuRequestsCores := cpuRequests.Value()

	// TODO need to add price filter only for Digital ocean
	// Build the SQL query
	q := `SELECT resource_id, description, slug, memory, memory_units, vcpus, disk, disk_units, price_monthly, price_hourly, region, cloud_provider
    	  FROM nodes
    	  WHERE cloud_provider = $1 AND region = $2 AND memory >= $3 AND vcpus >= $4 AND price_monthly >= 12
    	  ORDER BY price_hourly ASC`
	args := []interface{}{
		nf.CloudProvider,
		nf.Region,
		memRequestsMegaBytes,
		cpuRequestsCores,
	}
	// Execute the SQL query
	rows, err := apps.Pg.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the result into a NodesSlice
	nodes := hestia_autogen_bases.NodesSlice{}
	for rows.Next() {
		var node hestia_autogen_bases.Nodes
		err = rows.Scan(
			&node.ResourceID,
			&node.Description,
			&node.Slug,
			&node.Memory,
			&node.MemoryUnits,
			&node.Vcpus,
			&node.Disk,
			&node.DiskUnits,
			&node.PriceMonthly,
			&node.PriceHourly,
			&node.Region,
			&node.CloudProvider,
		)
		node.PriceHourly *= 1.1  // Add 10% to the price
		node.PriceMonthly *= 1.1 // Add 10% to the price

		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}
