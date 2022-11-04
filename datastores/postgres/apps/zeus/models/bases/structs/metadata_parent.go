package structs

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type ParentMetaData struct {
	autogen_bases.ChartSubcomponentParentClassTypes
	Metadata
}

func NewParentMetaData(parentClassTypeName string) ParentMetaData {
	cm := ParentMetaData{
		ChartSubcomponentParentClassTypes: autogen_bases.ChartSubcomponentParentClassTypes{ChartSubcomponentParentClassTypeName: parentClassTypeName},
		Metadata:                          NewMetadata(),
	}
	return cm
}

func (pm *ParentMetaData) SetParentClassTypeIDs(id int) {
	pm.ChartSubcomponentParentClassTypeID = id
	pm.SetMetadataParentClassTypeIDs(id)
}
