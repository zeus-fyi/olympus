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
	q.RawQuery = `INSERT INTO disks(size, units, description, type, price_monthly, price_hourly, region, cloud_provider) 
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, disk.Size, disk.Units, disk.Description, disk.Type, disk.PriceMonthly, disk.PriceHourly, disk.Region, disk.CloudProvider)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}
