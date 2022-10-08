package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
)

func ConvertKeyValueToChildValues(m map[string]string) common.ChildValuesSlice {
	cvs := common.ChildValuesSlice{}
	for k, v := range m {
		cv := ConvertKeyValueStringToChildValues(k, v)
		cvs = append(cvs, cv)
	}
	return cvs
}

func ConvertKeyValueStringToChildValues(k, v string) autogen_structs.ChartSubcomponentsChildValues {
	cv := autogen_structs.ChartSubcomponentsChildValues{
		ChartSubcomponentChildClassTypeID:              0,
		ChartSubcomponentChartPackageTemplateInjection: false,
		ChartSubcomponentKeyName:                       k,
		ChartSubcomponentValue:                         v,
	}
	return cv
}
