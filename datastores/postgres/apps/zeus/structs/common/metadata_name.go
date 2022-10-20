package common

type MetadataName struct {
	ChildClassSingleValue
}

func NewMetadataName() MetadataName {
	return MetadataName{NewChildClassSingleValue("name")}
}
