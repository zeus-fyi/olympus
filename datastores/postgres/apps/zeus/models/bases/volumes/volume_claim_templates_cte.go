package volumes

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (v *VolumeClaimTemplateGroup) GetVCTemplateGroupSubCTEs() sql_query_templates.SubCTEs {
	var combinedCTEs sql_query_templates.SubCTEs
	// TODO get parent
	for _, vt := range v.VolumeClaimTemplateSlice {
		combinedCTEs = sql_query_templates.AppendSubCteSlices(vt.GetVCTemplateSubCTEs(), combinedCTEs)
	}
	return combinedCTEs
}

func (v *VolumeClaimTemplate) GetVCTemplateSubCTEs() sql_query_templates.SubCTEs {
	vctMetadataSubCTEs := common.CreateMetadataSubCTEs(v.Metadata)
	vctSpecSubCTEs := v.Spec.GetVCTemplateSpecSubCTEs()
	return sql_query_templates.AppendSubCteSlices(vctMetadataSubCTEs, vctSpecSubCTEs)
}

func (v *VolumeClaimTemplateSpec) GetVCTemplateSpecSubCTEs() sql_query_templates.SubCTEs {
	vcStorageClassNameCTEs := common.CreateChildClassSingleValueSubCTEs(&v.StorageClassName)
	vcAccessModesCTEs := common.CreateChildClassMultiValueSubCTEs(&v.AccessModes)
	vcResourceRequestsCTEs := common.CreateChildClassMultiValueSubCTEs(&v.ResourceRequests)
	return sql_query_templates.AppendSubCteSlices(vcStorageClassNameCTEs, vcAccessModesCTEs, vcResourceRequestsCTEs)
}
