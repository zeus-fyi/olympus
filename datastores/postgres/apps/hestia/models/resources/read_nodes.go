package hestia_compute_resources

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"k8s.io/apimachinery/pkg/api/resource"
)

type NodeFilter struct {
	CloudProvider string                 `json:"cloudProvider"`
	Region        string                 `json:"region"`
	DiskType      string                 `json:"diskType,omitempty"`
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
	//0.5Gi = (1024*1024+512)*1 * 0.1 vCPU, adds this as overhead
	memRequestsMegaBytes := (memRequests.Value() + ((1024 * 1024 * 512) * 1)) / (1024 * 1024)
	cpuRequestsMilli := cpuRequests.MilliValue()
	cpuRequestsCores := float64(cpuRequestsMilli) / 1000

	switch strings.ToLower(nf.DiskType) {
	case "nvme":
		nf.DiskType = "nvme"
	default:
		nf.DiskType = "ssd"
	}
	// Build the SQL query
	q := `SELECT resource_id, description, slug, memory, memory_units, vcpus, disk, disk_units, price_monthly, price_hourly, region, cloud_provider, gpus, gpu_type
    	  FROM nodes
    	  WHERE memory >= $1 AND (vcpus + .1) >= $2 AND disk_type = $3
		  AND (
				(cloud_provider = 'do' AND price_monthly >= 12)
				OR
				(cloud_provider = 'gcp')
		      	OR 
				(cloud_provider = 'aws')
		      	OR 
				(cloud_provider = 'ovh')
			  )
		  AND (region = 'us-central1' OR region = 'nyc1' OR region = 'us-west-1' OR region = 'us-west-or-1')
    	AND price_monthly < 3000
		ORDER BY cloud_provider, price_hourly ASC;`
	args := []interface{}{
		memRequestsMegaBytes,
		cpuRequestsCores,
		nf.DiskType,
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
			&node.Gpus,
			&node.GpuType,
		)

		if err != nil {
			log.Err(err).Msg("Error scanning node")
			return nil, err
		}

		if node.CloudProvider == "gcp" {
			if strings.HasPrefix(node.Slug, "n2") {
				continue
			}
		}
		switch node.CloudProvider {
		case "do":
			node.PriceHourly *= 1.2  // Add 10% to the price
			node.PriceMonthly *= 1.2 // Add 10% to the price
		case "gcp":
			node.PriceHourly *= 1.30  // Add 40% to the price
			node.PriceMonthly *= 1.30 // Add 40% to the price
		case "aws":
			node.PriceHourly *= 1.30  // Add 40% to the price
			node.PriceMonthly *= 1.30 // Add 40% to the price
		case "ovh":
			node.PriceHourly *= 1.20  // Add 20% to the price
			node.PriceMonthly *= 1.20 // Add 20% to the price
		}
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
