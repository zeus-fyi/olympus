package hestia_compute_resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Resources represents a collection of nodes.
type Resources struct {
	Nodes hestia_autogen_bases.NodesSlice `json:"nodes"`
	// TODO: should use a map instead of a slice
	Disks hestia_autogen_bases.DisksSlice `json:"disks,omitempty"`
}

// RegionResourcesMap maps region names to their corresponding Resources.
type RegionResourcesMap map[string]Resources

// CloudProviderRegionsResourcesMap maps cloud provider names to their RegionResourcesMap,
// allowing for a nested mapping of providers to regions to resources.
type CloudProviderRegionsResourcesMap map[string]RegionResourcesMap

func SelectNodesV2(ctx context.Context, nf NodeFilter) (CloudProviderRegionsResourcesMap, error) {
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

	qa := ""
	qorg := ""

	args := []interface{}{
		memRequestsMegaBytes,
		cpuRequestsCores,
	}

	switch strings.ToLower(nf.DiskType) {
	case "nvme":
		nf.DiskType = "nvme"
		args = append(args, nf.DiskType)
		qa = fmt.Sprintf(" AND disk_type = $%d", len(args))
	case "ssd":
		nf.DiskType = "ssd"
		args = append(args, nf.DiskType)
		args = append(args, nf.DiskType)
		qa = fmt.Sprintf(" AND disk_type = $%d", len(args))
	default:
	}
	if nf.Ou.OrgID > 0 {
		args = append(args, nf.Ou.OrgID)
		qorg = fmt.Sprintf(" OR org_id = $%d", len(args))
	}

	// Build the SQL query
	q := `WITH user_auth_ctxs AS (
			SELECT
				cloud_provider,
				region
			FROM authorized_cluster_configs
			WHERE is_public = true ` + qorg + ` 
			GROUP BY cloud_provider, region
		  )
		  SELECT resource_id, description, slug, memory, memory_units, vcpus, disk, disk_units, price_monthly, price_hourly, n.region, n.cloud_provider, gpus, gpu_type
    	  FROM nodes n
		  JOIN user_auth_ctxs uac ON n.cloud_provider = uac.cloud_provider AND n.region = uac.region
    	  WHERE memory >= $1 AND (vcpus + .1) >= $2 ` + qa + `
			  AND (
					(n.cloud_provider = 'do' AND n.price_monthly >= 12)
					OR n.cloud_provider IN ('gcp', 'aws', 'ovh')
				  )
			  AND n.price_monthly < 3000
		ORDER BY cloud_provider, price_hourly ASC;`

	// Execute the SQL query
	rows, err := apps.Pg.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the result into a NodesSlice
	nodes := hestia_autogen_bases.NodesSlice{}
	cloudProviderRegionsResources := make(CloudProviderRegionsResourcesMap)

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
		di := hestia_autogen_bases.Disks{
			DiskUnits:     "Gi",
			Type:          node.DiskType,
			Region:        node.Region,
			CloudProvider: node.CloudProvider,
		}
		switch node.CloudProvider {
		case "do":
			di.PriceMonthly = 0.0137
			node.PriceHourly *= 1.00  // Add 10% to the price
			node.PriceMonthly *= 1.00 // Add 10% to the price
		case "gcp":
			di.PriceHourly = 0.02329
			node.PriceHourly *= 1.00  // Add 40% to the price
			node.PriceMonthly *= 1.00 // Add 40% to the price
		case "aws":
			di.PriceHourly = 0.01765
			node.PriceHourly *= 1.00  // Add 40% to the price
			node.PriceMonthly *= 1.00 // Add 40% to the price
		case "ovh":
			di.PriceHourly = 0.01643835616
			node.PriceHourly *= 1.00  // Add 20% to the price
			node.PriceMonthly *= 1.00 // Add 20% to the price
		}
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
		// Insert the node into the map
		if _, ok := cloudProviderRegionsResources[node.CloudProvider]; !ok {
			cloudProviderRegionsResources[node.CloudProvider] = make(RegionResourcesMap)
		}
		if _, ok := cloudProviderRegionsResources[node.CloudProvider][node.Region]; !ok {
			cloudProviderRegionsResources[node.CloudProvider][node.Region] = Resources{Nodes: make(hestia_autogen_bases.NodesSlice, 0)}
		}
		tmp := cloudProviderRegionsResources[node.CloudProvider][node.Region].Nodes
		tmp = append(tmp, node)
		cloudProviderRegionsResources[node.CloudProvider][node.Region] = Resources{Nodes: tmp}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cloudProviderRegionsResources, nil
}
