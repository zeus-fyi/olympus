package structs

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type ChildMetadata struct {
	autogen_bases.ChartSubcomponentChildClassTypes
	Metadata
}

func NewChildMetadata(childClassTypeName string, m Metadata) ChildMetadata {
	cm := ChildMetadata{
		ChartSubcomponentChildClassTypes: autogen_bases.ChartSubcomponentChildClassTypes{ChartSubcomponentChildClassTypeName: childClassTypeName},
		Metadata:                         m,
	}
	return cm
}

func (cm *ChildMetadata) SetChildClassTypeID(id int) {
	cm.Metadata.Name.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID = id
	cm.Metadata.Annotations.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID = id
	cm.Metadata.Labels.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID = id
	cm.ChartSubcomponentChildClassTypeID = id
}
