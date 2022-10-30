package structs

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

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

type Metadata struct {
	Name        ChildClassSingleValue
	Annotations ChildClassMultiValue
	Labels      ChildClassMultiValue
}

func NewMetadata() Metadata {
	m := Metadata{}
	m.Name = NewMetadataName()
	m.Annotations = NewMetadataAnnotations()
	m.Labels = NewMetadataLabels()
	return m
}

func (m *Metadata) SetMetadataParentClassTypeIDs(id int) {
	m.Name.ChartSubcomponentParentClassTypeID = id
	m.Annotations.ChartSubcomponentParentClassTypeID = id
	m.Labels.ChartSubcomponentParentClassTypeID = id
}

func (m *Metadata) HasName() bool {
	return len(m.Name.ChartSubcomponentsChildValues.ChartSubcomponentValue) > 0
}

func (m *Metadata) HasLabels() bool {
	return len(m.Labels.Values) > 0
}

func (m *Metadata) HasAnnotations() bool {
	return len(m.Annotations.Values) > 0
}
