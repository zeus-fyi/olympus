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
