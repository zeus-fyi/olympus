package fields

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// GenerateTags is where you add struct field tags. TODO make this tag map passed in as a param
func (f *Field) GenerateTags() map[string]string {
	tags := make(map[string]string)
	dbFieldName := f.DbFieldName()
	if len(dbFieldName) > 0 {
		jsonFieldTag := dbFieldName
		if strings.HasSuffix(dbFieldName, "_id") {
			jsonFieldTag = strings.TrimSuffix(dbFieldName, "_id")
			jsonFieldTag += "ID"
		}
		tags["db"] = dbFieldName
		tags["json"] = strcase.ToLowerCamel(jsonFieldTag)
	}
	return tags
}
