package fields

func (f *Field) GenerateTags() map[string]string {
	tags := make(map[string]string)
	dbFieldName := f.DbFieldName()
	if len(dbFieldName) > 0 {
		tags["db"] = dbFieldName
	}
	return tags
}
