package common

import autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"

type MetadataAnnotations struct {
	autogen_structs.ChartSubcomponentChildClassTypes
	AnnotationValues []autogen_structs.ChartSubcomponentsChildValues
}

func NewMetadataAnnotations() MetadataAnnotations {
	ma := MetadataAnnotations{
		autogen_structs.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: "annotations",
		}, []autogen_structs.ChartSubcomponentsChildValues{},
	}
	return ma
}

func (ma *MetadataAnnotations) AddAnnotations(annotations ...autogen_structs.ChartSubcomponentsChildValues) {
	if len(ma.AnnotationValues) <= 0 {
		ma.AnnotationValues = []autogen_structs.ChartSubcomponentsChildValues{}
	}
	for _, l := range annotations {
		ma.AnnotationValues = append(ma.AnnotationValues, l)
	}
}
