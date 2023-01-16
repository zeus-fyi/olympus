package validator_service_group

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

var ts chronos.Chronos

const Sn = "ValidatorServiceOrgGroup"

func InsertValidatorServiceOrgGroup(ctx context.Context, orgGroups hestia_autogen_bases.ValidatorServiceOrgGroupSlice, orgID int) (hestia_autogen_bases.ValidatorServiceOrgGroupSlice, error) {
	q := sql_query_templates.QueryParams{}
	cte := sql_query_templates.CTE{Name: "InsertValidatorServiceOrgGroup"}
	cte.SubCTEs = make([]sql_query_templates.SubCTE, len(orgGroups))
	cte.Params = []interface{}{}
	for i, orgGroup := range orgGroups {
		tmp := &orgGroup
		tmp.OrgID = &orgID
		tmp.Pubkey = strings_filter.AddHexPrefix(orgGroups[i].Pubkey)
		queryName := fmt.Sprintf("vsg_insert_%d", ts.UnixTimeStampNow())
		scte := sql_query_templates.NewSubInsertCTE(queryName)
		scte.TableName = tmp.GetTableName()
		scte.Columns = tmp.GetTableColumns()
		scte.Values = []apps.RowValues{tmp.GetRowValues(queryName)}
		cte.SubCTEs[i] = scte
		tmp.OrgID = nil
	}
	q.RawQuery = cte.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery, cte.Params...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return orgGroups, err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("Packages: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return orgGroups, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func InsertValidatorServiceOrgGroupCloudCtxNs(ctx context.Context, cloudCtxServiceGroup hestia_autogen_bases.ValidatorsServiceOrgGroupsCloudCtxNsSlice) (hestia_autogen_bases.ValidatorsServiceOrgGroupsCloudCtxNsSlice, error) {
	q := sql_query_templates.QueryParams{}
	cte := sql_query_templates.CTE{Name: "InsertValidatorServiceOrgGroupCloudCtxNs"}
	cte.SubCTEs = make([]sql_query_templates.SubCTE, len(cloudCtxServiceGroup))
	cte.Params = []interface{}{}
	for i, cloudGroup := range cloudCtxServiceGroup {
		tmp := &cloudGroup
		tmp.Pubkey = strings_filter.AddHexPrefix(cloudCtxServiceGroup[i].Pubkey)
		queryName := fmt.Sprintf("vsg_ns_insert_%d", ts.UnixTimeStampNow())
		scte := sql_query_templates.NewSubInsertCTE(queryName)
		scte.TableName = tmp.GetTableName()
		scte.Columns = tmp.GetTableColumns()
		scte.Values = []apps.RowValues{tmp.GetRowValues(queryName)}
		cte.SubCTEs[i] = scte
	}
	q.RawQuery = cte.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery, cte.Params...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return cloudCtxServiceGroup, err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("Packages: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return cloudCtxServiceGroup, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
