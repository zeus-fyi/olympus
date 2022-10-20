package common

type MetadataLabels struct {
	ChildClassAndValues
}

func NewMetadataLabels() MetadataLabels {
	ml := MetadataLabels{NewChildClassAndValues("labels")}
	return ml
}
