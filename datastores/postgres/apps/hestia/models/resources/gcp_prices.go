package hestia_compute_resources

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

var gcpPriceQuery = `
WITH selected_skus AS (
  SELECT *
  FROM gcp_services_skus
  WHERE (description LIKE '%' || $1 || '%') 
    AND usage_type = 'OnDemand'
    AND resource_family = 'Compute'
    AND (resource_group = 'CPU' OR resource_group='RAM' OR resource_group='N1Standard')
    AND service_regions::text LIKE '%' || $6 || '%'
), selected_sku_gpus AS (
  SELECT *
  FROM gcp_services_skus
  WHERE description LIKE '%' || $2 || '%'
    AND usage_type = 'OnDemand'
    AND resource_family = 'Compute'
    AND (resource_group = 'GPU')
    AND service_regions::text LIKE '%' || 'us-central1' || '%'
), gpu_pricing_data AS (
  SELECT
    resource_group,
    jsonb_array_elements(pricing_info) -> 'pricingExpression' ->> 'usageUnit' AS usage_unit,
    (jsonb_array_elements(pricing_info) -> 'pricingExpression' -> 'tieredRates' -> 0 -> 'unitPrice' ->> 'nanos')::float / 1000000000 AS unit_price
  FROM selected_sku_gpus
), pricing_data AS (
  SELECT
    resource_group,
	description,
    jsonb_array_elements(pricing_info) -> 'pricingExpression' ->> 'usageUnit' AS usage_unit,
    (jsonb_array_elements(pricing_info) -> 'pricingExpression' -> 'tieredRates' -> 0 -> 'unitPrice' ->> 'nanos')::float / 1000000000 AS unit_price
  FROM selected_skus
), gpu_cost AS (
  SELECT SUM(unit_price) * $3 AS total_gpu_cost
  FROM gpu_pricing_data
  WHERE usage_unit IN ('h', 'hour', 'hours') AND resource_group = 'GPU'
)
, cpu_cost AS (
  SELECT SUM(unit_price) * $4 AS total_cpu_cost
  FROM pricing_data
  WHERE usage_unit IN ('h', 'hour', 'hours') AND resource_group = 'CPU' OR (resource_group = 'N1Standard' AND description LIKE '%N1 Predefined Instance Core%')
)
, memory_cost AS (
  SELECT SUM(unit_price) * $5 AS total_memory_cost
  FROM pricing_data
  WHERE usage_unit IN ('GiBy.h', 'gibihours', 'gibihour') AND resource_group = 'RAM' OR (resource_group = 'N1Standard' AND description LIKE '%N1 Predefined Instance Ram%')
)
SELECT
  COALESCE((SELECT total_cpu_cost FROM cpu_cost), 0) AS total_cpu_cost,
  COALESCE((SELECT total_memory_cost FROM memory_cost), 0) AS total_memory_cost,
  COALESCE((SELECT total_gpu_cost FROM gpu_cost), 0) AS total_gpu_cost,
  COALESCE((SELECT total_cpu_cost FROM cpu_cost), 0) + COALESCE((SELECT total_memory_cost FROM memory_cost), 0) + COALESCE((SELECT total_gpu_cost FROM gpu_cost), 0) AS total_hourly_cost,
  (COALESCE((SELECT total_cpu_cost FROM cpu_cost), 0) + COALESCE((SELECT total_memory_cost FROM memory_cost), 0) + COALESCE((SELECT total_gpu_cost FROM gpu_cost), 0)) * 730 AS total_monthly_cost;
`

// SelectGcpPrices returns the total hourly and monthly cost of a GCP instance, TODO select by regions vs hardcoded
func SelectGcpPrices(ctx context.Context, name, gpuName string, gpuCount int, cpuCount, memSizeGB float64) (float64, float64, error) {
	var cpuUnitCost, memUnitCost, gpuUnitCost, totalHourlyCost, totalMonthlyCost float64
	region := "us-central1"
	err := apps.Pg.QueryRowWArgs(ctx, gcpPriceQuery, name, gpuName, gpuCount, cpuCount, memSizeGB, region).Scan(&cpuUnitCost, &memUnitCost, &gpuUnitCost, &totalHourlyCost, &totalMonthlyCost)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("SelectGcpPrices: Error")
		return totalHourlyCost, totalMonthlyCost, err
	}
	return totalHourlyCost, totalMonthlyCost, err
}
