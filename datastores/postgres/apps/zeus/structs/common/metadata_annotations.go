package common

type MetadataAnnotations struct {
	ChildClassAndValues
}

func NewMetadataAnnotations() MetadataAnnotations {
	ma := MetadataAnnotations{NewChildClassAndValues("annotations")}
	return ma
}
