package hestia_compute_resources

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/zeus/zeus/cluster_resources/nodes"
	"k8s.io/apimachinery/pkg/api/resource"
)

func SearchAndSelectNodes(ctx context.Context, nf nodes.NodeSearchParams) (hestia_autogen_bases.NodesSlice, error) {
	if nf.ResourceMinMax.Min.MemRequests == "" {
		nf.ResourceMinMax.Min.MemRequests = "0"
	}
	minMemRequests, err := resource.ParseQuantity(nf.ResourceMinMax.Min.MemRequests)
	if err != nil {
		return nil, err
	}

	if nf.ResourceMinMax.Min.CpuRequests == "" {
		nf.ResourceMinMax.Min.CpuRequests = "0"
	}
	cpuRequests, err := resource.ParseQuantity(nf.ResourceMinMax.Min.CpuRequests)
	if err != nil {
		return nil, err
	}
	// Convert to MegaBytes and vCores
	minMemRequestsMegaBytes := (minMemRequests.Value()) / (1024 * 1024)
	cpuRequestsMilli := cpuRequests.MilliValue()
	cpuRequestsCores := float64(cpuRequestsMilli) / 1000

	switch strings.ToLower(nf.DiskType) {
	case "nvme":
		nf.DiskType = "nvme"
	default:
		nf.DiskType = "ssd"
	}
	var cloudProviders []string
	var regionConditions []string
	for cp, regions := range nf.CloudProviderRegions {
		switch cp {
		case "aws":
			cloudProviders = append(cloudProviders, cp)
			for _, region := range regions {
				switch region {
				case "us-west-1":
					regionConditions = append(regionConditions, region)
				}
			}
		case "gcp":
			cloudProviders = append(cloudProviders, cp)
			for _, region := range regions {
				regionConditions = append(regionConditions, region)

			}
		case "do":
			cloudProviders = append(cloudProviders, cp)
			for _, region := range regions {
				regionConditions = append(regionConditions, region)
			}
		case "ovh":
			cloudProviders = append(cloudProviders, cp)
			for _, region := range regions {
				regionConditions = append(regionConditions, region)
			}
		}
	}
	addPriceSearch := ""
	args := []interface{}{
		minMemRequestsMegaBytes,
		cpuRequestsCores,
		nf.DiskType,
		pq.Array(cloudProviders),
		pq.Array(regionConditions),
	}
	if nf.ResourceMinMax.Min.MonthlyPrice != 0 {
		addPriceSearch += fmt.Sprintf(" AND price_monthly >= $%d", len(args)+1)
		args = append(args, nf.ResourceMinMax.Min.MonthlyPrice)
	}
	if nf.ResourceMinMax.Max.MonthlyPrice != 0 && nf.ResourceMinMax.Max.MonthlyPrice > nf.ResourceMinMax.Min.MonthlyPrice {
		addPriceSearch += fmt.Sprintf(" AND price_monthly <= $%d", len(args)+1)
		args = append(args, nf.ResourceMinMax.Max.MonthlyPrice)
	}
	if nf.ResourceMinMax.Max.MemRequests != "" {
		maxMemReq, merr := resource.ParseQuantity(nf.ResourceMinMax.Max.MemRequests)
		if merr != nil {
			return nil, merr
		}
		maxMemReqMegaBytesValue := (maxMemReq.Value()) / (1024 * 1024)
		if maxMemReqMegaBytesValue > 0 && maxMemReqMegaBytesValue > minMemRequestsMegaBytes {
			addPriceSearch += fmt.Sprintf(" AND memory <= $%d", len(args)+1)
			args = append(args, maxMemReqMegaBytesValue)
		}
	}

	if nf.ResourceMinMax.Max.CpuRequests != "" {
		maxMemReq, merr := resource.ParseQuantity(nf.ResourceMinMax.Max.CpuRequests)
		if merr != nil {
			return nil, merr
		}
		cpuRequestsMilliMax := maxMemReq.MilliValue()
		cpuRequestsCoresMax := float64(cpuRequestsMilliMax) / 1000
		if cpuRequestsCoresMax > 0 && cpuRequestsCoresMax > cpuRequestsCores {
			addPriceSearch += fmt.Sprintf(" AND vcpus <= $%d", len(args)+1)
			args = append(args, cpuRequestsCoresMax)
		}
	}

	// Build the SQL query
	q := fmt.Sprintf(`SELECT resource_id, description, slug, memory, memory_units, vcpus, disk, disk_units, price_monthly, price_hourly, region, cloud_provider, gpus, gpu_type
    	  FROM nodes
    	  WHERE memory >= $1 AND vcpus >= $2 AND disk_type = $3
		  AND (cloud_provider = ANY ($4::text[]))
		  AND (region = ANY ($5::text[])) %s
		ORDER BY cloud_provider, price_hourly ASC
		LIMIT 1000;`, addPriceSearch)

	// Execute the SQL query
	rows, err := apps.Pg.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Parse the result into a NodesSlice
	nodesSlice := hestia_autogen_bases.NodesSlice{}
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
		node.PriceHourly = roundFloat(node.PriceHourly, 2)
		node.PriceMonthly = roundFloat(node.PriceMonthly, 2)
		if err != nil {
			return nil, err
		}
		nodesSlice = append(nodesSlice, node)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nodesSlice, nil
}
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
