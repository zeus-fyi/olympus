package hestia_compute_resources

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertDisk(ctx context.Context, disk hestia_autogen_bases.Disks) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `With cte_insert_resource AS (
					INSERT INTO resources(type)
					VALUES ('disk')
					RETURNING resource_id
				  )				
				  INSERT INTO disks(resource_id, disk_size, disk_units, description, type, price_monthly, price_hourly, region, cloud_provider) 
				  VALUES ((SELECT resource_id FROM cte_insert_resource), $1, $2, $3, $4, $5, $6, $7, $8)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, disk.DiskSize, disk.DiskUnits, disk.Description, disk.Type, disk.PriceMonthly, disk.PriceHourly, disk.Region, disk.CloudProvider)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func GetDiskTypesAWS(region string) []hestia_autogen_bases.Disks {
	gbInGi := 107.374

	resourceStrID := ""
	if region == "us-west-1" {
		gbInGi *= 0.1
		resourceStrID = "1707332103514001000"
	} else if region == "us-east-2" {
		gbInGi *= 0.08
		resourceStrID = "1707332235297634000"
	} else if region == "us-east-1" {
		gbInGi *= 0.08
		resourceStrID = "1707486802147690000"
	}
	disk := hestia_autogen_bases.Disks{
		ResourceStrID: resourceStrID,
		PriceMonthly:  gbInGi,
		PriceHourly:   gbInGi / 730,
		Region:        region,
		CloudProvider: "aws",
		Description:   "EBS gp3 Block Storage SSD",
		Type:          "ssd",
		SubType:       "gp3",
		DiskSize:      100,
		DiskUnits:     "Gi",
	}
	gbInGi = 107.374
	if region == "us-west-1" {
		gbInGi *= 0.12
		resourceStrID = "1683860918169422000"
	} else if region == "us-east-2" {
		gbInGi *= 0.1
		resourceStrID = "1707332260863081000"
	} else if region == "us-east-1" {
		gbInGi *= 0.1
		resourceStrID = "1707486750716592000"
	}
	disk2 := hestia_autogen_bases.Disks{
		ResourceStrID: resourceStrID,
		PriceMonthly:  gbInGi,
		PriceHourly:   gbInGi / 730,
		Region:        region,
		CloudProvider: "aws",
		Description:   "EBS gp2 Block Storage SSD",
		Type:          "ssd",
		SubType:       "gp2",
		DiskSize:      100,
		DiskUnits:     "Gi",
	}
	return []hestia_autogen_bases.Disks{
		disk, disk2,
	}
}
