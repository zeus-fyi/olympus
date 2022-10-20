package common

type Metadata struct {
	Name        MetadataName
	Annotations MetadataAnnotations
	Labels      MetadataLabels
}

func NewMetadata() Metadata {
	m := Metadata{}
	m.Name = NewMetadataName()
	m.Annotations = NewMetadataAnnotations()
	m.Labels = NewMetadataLabels()
	return m
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
