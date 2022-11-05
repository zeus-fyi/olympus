package volumes

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (v *VolumeClaimTemplateGroup) GetVCTemplateGroupSubCTEs(c *charts.Chart) sql_query_templates.SubCTEs {
	ts := chronos.Chronos{}
	if v.ChartSubcomponentParentClassTypeID == 0 {
		v.SetParentIDs(ts.UnixTimeStampNow())
	}
	v.SetNewChildIDs()

	var combinedCTEs sql_query_templates.SubCTEs
	pcSubCTEs := common.CreateParentClassTypeSubCTE(c, &v.ParentClass.ChartSubcomponentParentClassTypes)
	pkgSubCTE := common.AddParentClassToChartPackage(c, v.ChartSubcomponentParentClassTypeID)
	combinedCTEs = sql_query_templates.AppendSubCteSlices(pcSubCTEs)
	for _, vt := range v.VolumeClaimTemplateSlice {
		combinedCTEs = sql_query_templates.AppendSubCteSlices(vt.GetVCTemplateSubCTEs(), combinedCTEs)
	}
	return sql_query_templates.AppendSubCteSlices(combinedCTEs, []sql_query_templates.SubCTE{pkgSubCTE})
}

func (v *VolumeClaimTemplate) GetVCTemplateSubCTEs() sql_query_templates.SubCTEs {
	vctMetadataSubCTEs := common.CreateBaseMetadataSubCTEs(v.Metadata)
	vctSpecSubCTEs := v.Spec.GetVCTemplateSpecSubCTEs()
	return sql_query_templates.AppendSubCteSlices(vctMetadataSubCTEs, vctSpecSubCTEs)
}

func (v *VolumeClaimTemplateSpec) GetVCTemplateSpecSubCTEs() sql_query_templates.SubCTEs {
	vcStorageClassNameCTEs := common.CreateChildClassSingleValueSubCTEs(&v.StorageClassName)
	vcAccessModesCTEs := common.CreateChildClassMultiValueSubCTEs(&v.AccessModes)
	vcResourceRequestsCTEs := common.CreateChildClassMultiValueSubCTEs(&v.ResourceRequests)
	return sql_query_templates.AppendSubCteSlices(vcStorageClassNameCTEs, vcAccessModesCTEs, vcResourceRequestsCTEs)
}
