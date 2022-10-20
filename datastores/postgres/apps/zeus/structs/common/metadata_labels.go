package common

import autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"

type MetadataLabels struct {
	autogen_structs.ChartSubcomponentChildClassTypes
	LabelValues []autogen_structs.ChartSubcomponentsChildValues
}

func NewMetadataLabels() MetadataLabels {
	ml := MetadataLabels{
		autogen_structs.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: "labels",
		}, []autogen_structs.ChartSubcomponentsChildValues{},
	}
	return ml
}

func (ml *MetadataLabels) AddLabels(labels ...autogen_structs.ChartSubcomponentsChildValues) {
	if len(ml.LabelValues) <= 0 {
		ml.LabelValues = []autogen_structs.ChartSubcomponentsChildValues{}
	}
	for _, l := range labels {
		ml.LabelValues = append(ml.LabelValues, l)
	}
}
