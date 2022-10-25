package common_conversions

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs/common"
)

func ConvertKeyValueToChildValues(m map[string]string) common.ChildValuesSlice {
	cvs := common.ChildValuesSlice{}
	for k, v := range m {
		cv := ConvertKeyValueStringToChildValues(k, v)
		cvs = append(cvs, cv)
	}
	return cvs
}

func ConvertKeyValueStringToChildValues(k, v string) autogen_bases.ChartSubcomponentsChildValues {
	cv := autogen_bases.ChartSubcomponentsChildValues{
		ChartSubcomponentChildClassTypeID:              0,
		ChartSubcomponentChartPackageTemplateInjection: false,
		ChartSubcomponentKeyName:                       k,
		ChartSubcomponentValue:                         v,
	}
	return cv
}
