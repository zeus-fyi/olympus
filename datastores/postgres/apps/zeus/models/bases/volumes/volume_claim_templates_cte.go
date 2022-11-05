package volumes

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (v *VolumeClaimTemplateGroup) GetVCTemplateGroupSubCTEs(c *charts.Chart) sql_query_templates.SubCTEs {
	ts := chronos.Chronos{}
	var combinedCTEs sql_query_templates.SubCTEs

	for _, vt := range v.VolumeClaimTemplateSlice {
		if vt.ChartSubcomponentParentClassTypeID == 0 {
			vt.ChartSubcomponentParentClassTypeID = ts.UnixTimeStampNow()
		}
		vt.ChartComponentResourceID = c.ChartComponentResourceID
		vt.Metadata.ChartComponentResourceID = c.ChartComponentResourceID
		vt.Metadata.SetParentClassTypeIDs(vt.ChartSubcomponentParentClassTypeID)

		vt.Spec.SetParentClassTypeID(vt.ChartSubcomponentParentClassTypeID)
		vctMetadataSubCTEs := common.CreateParentMetadataSubCTEs(c, vt.Metadata)
		combinedCTEs = sql_query_templates.AppendSubCteSlices(vctMetadataSubCTEs, vt.GetVCTemplateSubCTEs(), combinedCTEs)
	}
	return combinedCTEs
}

func (v *VolumeClaimTemplate) GetVCTemplateSubCTEs() sql_query_templates.SubCTEs {
	vctSpecSubCTEs := v.Spec.GetVCTemplateSpecSubCTEs()
	return sql_query_templates.AppendSubCteSlices(vctSpecSubCTEs)
}

func (v *VolumeClaimTemplateSpec) GetVCTemplateSpecSubCTEs() sql_query_templates.SubCTEs {
	vcStorageClassNameCTEs := common.CreateChildClassSingleValueSubCTEs(&v.StorageClassName)
	vcAccessModesCTEs := common.CreateChildClassMultiValueSubCTEs(&v.AccessModes)
	vcResourceRequestsCTEs := common.CreateChildClassMultiValueSubCTEs(&v.ResourceRequests)
	return sql_query_templates.AppendSubCteSlices(vcStorageClassNameCTEs, vcAccessModesCTEs, vcResourceRequestsCTEs)
}
